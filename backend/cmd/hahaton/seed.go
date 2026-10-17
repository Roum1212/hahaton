package hahaton

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_cooperation_command "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation/command"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_event_chat "github.com/Roum1212/hahaton/internal/domain/aggregate/event-chat"
	aggregate_event_command "github.com/Roum1212/hahaton/internal/domain/aggregate/event/command"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	aggregate_participation_command "github.com/Roum1212/hahaton/internal/domain/aggregate/participation/command"
	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
	aggregate_reward_command "github.com/Roum1212/hahaton/internal/domain/aggregate/reward/command"
	model_actor_kind "github.com/Roum1212/hahaton/internal/domain/model/actor-kind"
	model_cooperation_change "github.com/Roum1212/hahaton/internal/domain/model/cooperation-change"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_event_requirements "github.com/Roum1212/hahaton/internal/domain/model/event-requirements"
	model_event_schedule "github.com/Roum1212/hahaton/internal/domain/model/event-schedule"
	model_event_title "github.com/Roum1212/hahaton/internal/domain/model/event-title"
	model_geo_point "github.com/Roum1212/hahaton/internal/domain/model/geo-point"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_operation_id "github.com/Roum1212/hahaton/internal/domain/model/operation-id"
	model_participation_change "github.com/Roum1212/hahaton/internal/domain/model/participation-change"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	model_participation_status "github.com/Roum1212/hahaton/internal/domain/model/participation-status"
	model_reward_id "github.com/Roum1212/hahaton/internal/domain/model/reward-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	model_volunteer_limit "github.com/Roum1212/hahaton/internal/domain/model/volunteer-limit"
	client_encryption "github.com/Roum1212/hahaton/internal/infra/client/encryption"
	client_postgresql "github.com/Roum1212/hahaton/internal/infra/client/postgresql"
	repository_cooperation_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/cooperation/postgresql"
	repository_coordinator_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/coordinator/postgresql"
	repository_event_chat_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/event-chat/postgresql"
	repository_event_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/event/postgresql"
	repository_participation_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/participation/postgresql"
	repository_reward_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/reward/postgresql"
	repository_volunteer_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/volunteer/postgresql"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

const (
	demoVolunteerLimit = 10
	demoRewardPrice    = 20
	demoRewardQuantity = 10
	demoLatitude       = 55.1
	demoLongitude      = 36.6
)

type seed struct {
	database       client_postgresql.Client
	coordinators   repository_coordinator_postgresql.Coordinators
	volunteers     repository_volunteer_postgresql.Volunteers
	cooperations   repository_cooperation_postgresql.Cooperations
	events         repository_event_postgresql.Events
	chats          *repository_event_chat_postgresql.EventChats
	participations repository_participation_postgresql.Participations

	rewards repository_reward_postgresql.Rewards
}

func Seed(
	ctx context.Context,
	pool *pgxpool.Pool,
	coordinatorID model_user_id.UserID,
	volunteerID model_user_id.UserID,
	chatKey []byte,
) error {
	if coordinatorID == volunteerID {
		return errors.New("demo coordinator and volunteer must be different participants")
	}

	x, err := newSeed(pool, chatKey)
	if err != nil {
		return err
	}

	if _, err := x.coordinators.Load(ctx, coordinatorID); err != nil {
		return err
	}

	if _, err := x.volunteers.Load(ctx, volunteerID); err != nil {
		return err
	}

	return x.database.InTx(ctx, func(ctx context.Context) error {
		now := time.Now().UTC().Truncate(time.Minute)
		if err := x.seedCooperation(ctx, coordinatorID, volunteerID, now); err != nil {
			return err
		}

		if err := x.seedEvents(ctx, coordinatorID, volunteerID, now); err != nil {
			return err
		}

		return x.seedRewards(ctx, coordinatorID, now)
	})
}

func (x seed) recordChange(
	ctx context.Context,
	value *aggregate_participation.Participation,
	previous utils_value.Optional[model_participation_status.ParticipationStatus],
	actorID model_user_id.UserID,
	actorKind model_actor_kind.ActorKind,
) error {
	id, err := model_operation_id.GenerateOperationID()
	if err != nil {
		return err
	}

	change, err := model_participation_change.NewParticipationChange(
		id,
		value.GetParticipationID(),
		previous,
		value.GetStatus(),
		actorID,
		actorKind,
		"",
		"",
		0,
		value.GetUpdatedAt(),
	)
	if err != nil {
		return err
	}

	return x.participations.AppendChange(ctx, change)
}

func (x seed) seedCooperation(
	ctx context.Context,
	coordinatorID model_user_id.UserID,
	volunteerID model_user_id.UserID,
	now time.Time,
) error {
	cooperation, err := x.cooperations.Load(ctx, coordinatorID, volunteerID)
	exists := true

	if err != nil {
		switch {
		case errors.Is(err, aggregate_cooperation.ErrNotFound):
			cooperation = aggregate_cooperation.NewCooperation(coordinatorID, volunteerID)
			exists = false
		default:
			return err
		}
	}

	previousStatus := cooperation.GetStatus()
	if err := cooperation.AcceptInvitation(aggregate_cooperation_command.NewAcceptInvitation(now)); err != nil {
		return err
	}

	recordChange := func(ctx context.Context) error {
		if cooperation.GetStatus() == previousStatus {
			return nil
		}

		change, err := model_cooperation_change.NewCooperationChange(
			coordinatorID,
			volunteerID,
			cooperation.GetStatus(),
			coordinatorID,
			cooperation.GetRetryAfter(),
			now,
		)
		if err != nil {
			return err
		}

		return x.cooperations.AppendChange(ctx, change)
	}
	if exists {
		if err := x.cooperations.Save(ctx, cooperation); err != nil {
			return err
		}

		return recordChange(ctx)
	}

	return x.cooperations.Create(ctx, cooperation, recordChange)
}

func (x seed) seedEvent(
	ctx context.Context,
	coordinatorID model_user_id.UserID,
	volunteerID model_user_id.UserID,
	titleText string,
	index int,
	now time.Time,
) error {
	point, err := model_geo_point.NewGeoPoint(demoLatitude, demoLongitude)
	if err != nil {
		return err
	}

	requirements, err := model_event_requirements.NewEventRequirements(nil)
	if err != nil {
		return err
	}

	limit, err := model_volunteer_limit.NewVolunteerLimit(demoVolunteerLimit)
	if err != nil {
		return err
	}

	title, err := model_event_title.NewEventTitle(titleText)
	if err != nil {
		return err
	}

	eventID, err := model_event_id.GenerateEventID()
	if err != nil {
		return err
	}

	start := now.Add(time.Duration(index+1) * 24 * time.Hour)

	createdAt := now
	if index == 0 {
		start = now.Add(-2 * time.Hour)
		createdAt = start.Add(-time.Hour)
	}

	schedule, err := model_event_schedule.NewEventSchedule(start, start.Add(time.Hour))
	if err != nil {
		return err
	}

	event := aggregate_event.NewEvent(eventID, coordinatorID)
	if err := event.Create(aggregate_event_command.NewCreate(
		title,
		"Учебные данные для защиты",
		model_image_path.ImagePath{},
		"Учебная точка на карте",
		point,
		schedule,
		requirements,
		limit,
		createdAt,
	)); err != nil {
		return err
	}

	if err := x.events.Create(ctx, event, func(ctx context.Context) error {
		return x.chats.Create(ctx, aggregate_event_chat.NewEventChat(eventID))
	}); err != nil {
		return err
	}

	if index != 0 {
		return nil
	}

	participationID, err := model_participation_id.GenerateParticipationID()
	if err != nil {
		return err
	}

	value := aggregate_participation.NewParticipation(participationID, eventID, coordinatorID, volunteerID)
	if err := value.Request(aggregate_participation_command.NewRequest(createdAt)); err != nil {
		return err
	}

	if err := x.participations.Create(ctx, value); err != nil {
		return err
	}

	if err := x.recordChange(
		ctx,
		value,
		utils_value.Optional[model_participation_status.ParticipationStatus]{},
		volunteerID,
		model_actor_kind.Volunteer,
	); err != nil {
		return err
	}

	previous := utils_value.NewOptional(value.GetStatus())
	if err := value.Approve(aggregate_participation_command.NewApprove(0, createdAt)); err != nil {
		return err
	}

	if err := x.participations.Save(ctx, value); err != nil {
		return err
	}

	return x.recordChange(ctx, value, previous, coordinatorID, model_actor_kind.Coordinator)
}

func (x seed) seedEvents(
	ctx context.Context,
	coordinatorID model_user_id.UserID,
	volunteerID model_user_id.UserID,
	now time.Time,
) error {
	values, err := x.events.List(ctx)
	if err != nil {
		return err
	}

	existing := make(map[string]bool)

	for _, event := range values {
		if event.GetCoordinatorID() == coordinatorID {
			existing[event.GetTitle().String()] = true
		}
	}

	for index, title := range []string{
		"[demo] Подготовленная акция",
		"[demo] Помощь библиотеке",
		"[demo] Городской субботник",
	} {
		if existing[title] {
			continue
		}

		if err := x.seedEvent(ctx, coordinatorID, volunteerID, title, index, now); err != nil {
			return err
		}
	}

	return nil
}

func (x seed) seedRewards(ctx context.Context, coordinatorID model_user_id.UserID, now time.Time) error {
	values, err := x.rewards.List(ctx)
	if err != nil {
		return err
	}

	existing := make(map[string]bool)

	for _, reward := range values {
		if reward.GetCoordinatorID() == coordinatorID {
			existing[reward.GetTitle()] = true
		}
	}

	for _, title := range []string{
		"[demo] Значок",
		"[demo] Открытка",
		"[demo] Благодарность",
	} {
		if existing[title] {
			continue
		}

		id, err := model_reward_id.GenerateRewardID()
		if err != nil {
			return err
		}

		reward := aggregate_reward.NewReward(id, coordinatorID)
		if err := reward.Create(aggregate_reward_command.NewCreate(
			title,
			"Учебное поощрение без реальной выдачи",
			model_image_path.ImagePath{},
			demoRewardPrice,
			demoRewardQuantity,
			now,
		)); err != nil {
			return err
		}

		if err := x.rewards.Create(ctx, reward); err != nil {
			return err
		}
	}

	return nil
}

func newSeed(pool *pgxpool.Pool, chatKey []byte) (seed, error) {
	database := client_postgresql.NewClient(pool)

	crypto, err := client_encryption.NewClient(chatKey)
	if err != nil {
		return seed{}, err
	}

	return seed{
		database:       database,
		coordinators:   repository_coordinator_postgresql.NewCoordinators(database),
		volunteers:     repository_volunteer_postgresql.NewVolunteers(database),
		cooperations:   repository_cooperation_postgresql.NewCooperations(database),
		events:         repository_event_postgresql.NewEvents(database),
		chats:          repository_event_chat_postgresql.NewEventChats(database, crypto),
		participations: repository_participation_postgresql.NewParticipations(database),

		rewards: repository_reward_postgresql.NewRewards(database),
	}, nil
}
