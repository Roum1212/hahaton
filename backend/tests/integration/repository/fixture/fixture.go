//go:build integration

package repository_fixture

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_cooperation_invitation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation"
	aggregate_cooperation_invitation_command "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation/command"
	aggregate_cooperation_command "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation/command"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_coordinator_command "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator/command"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_event_command "github.com/Roum1212/hahaton/internal/domain/aggregate/event/command"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	aggregate_participation_command "github.com/Roum1212/hahaton/internal/domain/aggregate/participation/command"
	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
	aggregate_reward_order_command "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order/command"
	aggregate_reward_command "github.com/Roum1212/hahaton/internal/domain/aggregate/reward/command"
	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
	aggregate_session_command "github.com/Roum1212/hahaton/internal/domain/aggregate/session/command"
	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	aggregate_volunteer_command "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer/command"
	aggregate_wallet "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet"
	aggregate_wallet_command "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet/command"
	model_date_of_birth "github.com/Roum1212/hahaton/internal/domain/model/date-of-birth"
	model_email_address "github.com/Roum1212/hahaton/internal/domain/model/email-address"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_event_requirement "github.com/Roum1212/hahaton/internal/domain/model/event-requirement"
	model_event_requirements "github.com/Roum1212/hahaton/internal/domain/model/event-requirements"
	model_event_schedule "github.com/Roum1212/hahaton/internal/domain/model/event-schedule"
	model_event_title "github.com/Roum1212/hahaton/internal/domain/model/event-title"
	model_full_name "github.com/Roum1212/hahaton/internal/domain/model/full-name"
	model_geo_point "github.com/Roum1212/hahaton/internal/domain/model/geo-point"
	model_identity "github.com/Roum1212/hahaton/internal/domain/model/identity"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_invitation_id "github.com/Roum1212/hahaton/internal/domain/model/invitation-id"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	model_phone_number "github.com/Roum1212/hahaton/internal/domain/model/phone-number"
	model_reward_id "github.com/Roum1212/hahaton/internal/domain/model/reward-id"
	model_reward_order_id "github.com/Roum1212/hahaton/internal/domain/model/reward-order-id"
	model_session_id "github.com/Roum1212/hahaton/internal/domain/model/session-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	model_volunteer_limit "github.com/Roum1212/hahaton/internal/domain/model/volunteer-limit"
	client_postgresql "github.com/Roum1212/hahaton/internal/infra/client/postgresql"
	repository_coordinator "github.com/Roum1212/hahaton/internal/infra/repository/coordinator/postgresql"
	repository_event "github.com/Roum1212/hahaton/internal/infra/repository/event/postgresql"
	repository_reward "github.com/Roum1212/hahaton/internal/infra/repository/reward/postgresql"
	repository_volunteer "github.com/Roum1212/hahaton/internal/infra/repository/volunteer/postgresql"
	types_uuid_v7 "github.com/Roum1212/hahaton/internal/pkg/types/uuid-v7"
)

func CreateCoordinator(t *testing.T, database client_postgresql.Client) *aggregate_coordinator.Coordinator {
	t.Helper()

	value := NewCoordinator(t, CreateVolunteer(t, database).GetUserID())
	require.NoError(t, repository_coordinator.NewCoordinators(database).Create(t.Context(), value))

	return value
}

func CreateEvent(t *testing.T, database client_postgresql.Client) *aggregate_event.Event {
	t.Helper()

	value := NewEvent(t, CreateCoordinator(t, database).GetUserID())
	require.NoError(
		t,
		repository_event.NewEvents(database).Create(t.Context(), value, func(ctx context.Context) error {
			return nil
		}),
	)

	return value
}

func CreateReward(t *testing.T, database client_postgresql.Client) *aggregate_reward.Reward {
	t.Helper()

	value := NewReward(t, CreateCoordinator(t, database).GetUserID())
	require.NoError(t, repository_reward.NewRewards(database).Create(t.Context(), value))

	return value
}

func CreateVolunteer(t *testing.T, database client_postgresql.Client) *aggregate_volunteer.Volunteer {
	t.Helper()

	value := NewVolunteer(t)
	require.NoError(
		t,
		repository_volunteer.NewVolunteers(database).Create(
			t.Context(),
			value,
			func(ctx context.Context) error {
				return nil
			},
		),
	)

	return value
}

func ImagePath(t *testing.T, owner model_user_id.UserID) model_image_path.ImagePath {
	t.Helper()

	imageID, err := types_uuid_v7.GenerateUUIDv7()
	require.NoError(t, err)
	path, err := model_image_path.NewImagePathFromFields(owner, imageID, ".png")
	require.NoError(t, err)

	return path
}

func NewCooperation(t *testing.T, owner, volunteer model_user_id.UserID) *aggregate_cooperation.Cooperation {
	t.Helper()

	value := aggregate_cooperation.NewCooperation(owner, volunteer)
	require.NoError(t, value.Request(aggregate_cooperation_command.NewRequest(Now())))

	return value
}

func NewCoordinator(t *testing.T, id model_user_id.UserID) *aggregate_coordinator.Coordinator {
	t.Helper()

	value := aggregate_coordinator.NewCoordinator(id)
	require.NoError(
		t,
		value.Create(aggregate_coordinator_command.NewCreate("Координатор", ImagePath(t, id), 60, Now())),
	)

	return value
}

func NewEvent(t *testing.T, owner model_user_id.UserID) *aggregate_event.Event {
	t.Helper()

	id, err := model_event_id.GenerateEventID()
	require.NoError(t, err)
	title, err := model_event_title.NewEventTitle("Доброе дело")
	require.NoError(t, err)
	point, err := model_geo_point.NewGeoPoint(55.1, 36.6)
	require.NoError(t, err)
	schedule, err := model_event_schedule.NewEventSchedule(Now().Add(time.Hour), Now().Add(4*time.Hour))
	require.NoError(t, err)
	requirements, err := model_event_requirements.NewEventRequirements([]model_event_requirement.Requirement{
		model_event_requirement.Age16Plus,
		model_event_requirement.PreliminaryTraining,
	})
	require.NoError(t, err)
	limit, err := model_volunteer_limit.NewVolunteerLimit(3)
	require.NoError(t, err)

	value := aggregate_event.NewEvent(id, owner)
	require.NoError(t, value.Create(aggregate_event_command.NewCreate(
		title, "Описание", ImagePath(t, owner), "Обнинск", point, schedule, requirements, limit, Now(),
	)))

	return value
}

func NewInvitation(t *testing.T, owner model_user_id.UserID) *aggregate_cooperation_invitation.CooperationInvitation {
	t.Helper()

	id, err := model_invitation_id.GenerateInvitationID()
	require.NoError(t, err)

	value := aggregate_cooperation_invitation.NewCooperationInvitation(id, owner)
	value.Create(
		aggregate_cooperation_invitation_command.NewCreate("invitation-"+id.String(), Now().Add(time.Hour), Now()),
	)

	return value
}

func NewOrder(
	t *testing.T,
	reward *aggregate_reward.Reward,
	user model_user_id.UserID,
) *aggregate_reward_order.RewardOrder {
	t.Helper()

	id, err := model_reward_order_id.GenerateRewardOrderID()
	require.NoError(t, err)

	value := aggregate_reward_order.NewRewardOrder(id, reward.GetRewardID(), reward.GetCoordinatorID(), user)
	require.NoError(t, value.Create(aggregate_reward_order_command.NewCreate("Сувенир", 20, id.String(), Now())))

	return value
}

func NewParticipation(
	t *testing.T,
	event *aggregate_event.Event,
	user model_user_id.UserID,
) *aggregate_participation.Participation {
	t.Helper()

	id, err := model_participation_id.GenerateParticipationID()
	require.NoError(t, err)

	value := aggregate_participation.NewParticipation(id, event.GetEventID(), event.GetCoordinatorID(), user)
	require.NoError(t, value.Request(aggregate_participation_command.NewRequest(Now())))

	return value
}

func NewReward(t *testing.T, owner model_user_id.UserID) *aggregate_reward.Reward {
	t.Helper()

	id, err := model_reward_id.GenerateRewardID()
	require.NoError(t, err)

	value := aggregate_reward.NewReward(id, owner)
	require.NoError(t, value.Create(aggregate_reward_command.NewCreate(
		"Сувенир", "Описание", ImagePath(t, owner), 20, 3, Now(),
	)))

	return value
}

func NewSession(t *testing.T) *aggregate_session.Session {
	t.Helper()

	id, err := model_session_id.GenerateSessionID()
	require.NoError(t, err)

	value := aggregate_session.NewSession(id)
	require.NoError(t, value.Open(aggregate_session_command.NewOpen(Now(), Now().Add(time.Hour))))

	return value
}

func NewVolunteer(t *testing.T) *aggregate_volunteer.Volunteer {
	t.Helper()

	id, err := model_user_id.GenerateUserID()
	require.NoError(t, err)
	identity, err := model_identity.NewIdentity("https://id.flida.test", id.String())
	require.NoError(t, err)
	name, err := model_full_name.NewFullName("Иван", "Иванов", "Иванович")
	require.NoError(t, err)
	date, err := model_date_of_birth.NewDateOfBirth(2000, time.February, 29)
	require.NoError(t, err)
	email, err := model_email_address.NewEmailAddress(id.String() + "@example.com")
	require.NoError(t, err)
	phone, err := model_phone_number.NewPhoneNumber("+79990000001")
	require.NoError(t, err)

	value := aggregate_volunteer.NewVolunteer(id)
	require.NoError(
		t,
		value.Register(aggregate_volunteer_command.NewRegister(identity, name, date, email, phone, Now())),
	)

	value.Update(aggregate_volunteer_command.NewUpdate(name, date, ImagePath(t, id), Now()))

	return value
}

func NewWallet(t *testing.T, owner model_user_id.UserID) *aggregate_wallet.Wallet {
	t.Helper()

	value := aggregate_wallet.NewWallet(owner)
	require.NoError(t, value.Open(aggregate_wallet_command.NewOpen(Now())))

	return value
}

func Now() time.Time {
	return time.Date(2026, time.October, 1, 12, 0, 0, 0, time.Local)
}

func Reset(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	_, err := pool.Exec(t.Context(), `TRUNCATE TABLE
 sessions, volunteers, wallets, coordinators, cooperations, cooperation_invitations,
 events, participations, participation_changes, volunteer_operations, rewards,
 reward_orders, event_chats, chat_messages, cities, cooperation_changes RESTART IDENTITY CASCADE`)
	require.NoError(t, err)
}
