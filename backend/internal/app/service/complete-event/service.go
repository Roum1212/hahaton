package service_complete_event

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_event_command "github.com/Roum1212/hahaton/internal/domain/aggregate/event/command"
	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	aggregate_participation_command "github.com/Roum1212/hahaton/internal/domain/aggregate/participation/command"
	aggregate_wallet "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet"
	aggregate_wallet_command "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet/command"
	model_actor_kind "github.com/Roum1212/hahaton/internal/domain/model/actor-kind"
	model_event_result "github.com/Roum1212/hahaton/internal/domain/model/event-result"
	model_operation_id "github.com/Roum1212/hahaton/internal/domain/model/operation-id"
	model_participation_change "github.com/Roum1212/hahaton/internal/domain/model/participation-change"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	model_participation_status "github.com/Roum1212/hahaton/internal/domain/model/participation-status"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

var (
	ErrAlreadyExists    = errors.New("already exists")
	ErrInvalidRequest   = errors.New("invalid request")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	Service interface {
		CompleteEvent(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		events         aggregate_event.Repository
		coordinators   aggregate_coordinator.Repository
		participations aggregate_participation.Repository

		wallets aggregate_wallet.Repository
	}
)

func (x service) CompleteEvent(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	if r.IsForced() {
		event, err := x.events.Load(ctx, r.GetEventID())
		if err != nil {
			return Response{}, err
		}

		r.userID = event.GetCoordinatorID()
	}

	if _, err := x.coordinators.Load(ctx, r.GetUserID()); err != nil {
		switch {
		case errors.Is(err, aggregate_coordinator.ErrNotFound):
			return Response{}, ErrPermissionDenied
		default:
			return Response{}, err
		}
	}

	fingerprint := resultFingerprint(r.GetResults())

	var completed *aggregate_event.Event

	if err := x.events.Modify(ctx, r.GetEventID(), func(
		ctx context.Context,
		event *aggregate_event.Event,
	) error {
		events, err := x.events.List(ctx)
		if err != nil {
			return err
		}

		for _, previous := range events {
			if previous.GetCoordinatorID() != r.GetUserID() || previous.GetCompletionRequestID() != r.GetRequestID() {
				continue
			}

			if previous.GetEventID() != r.GetEventID() || previous.GetResultFingerprint() != fingerprint {
				return ErrAlreadyExists
			}

			completed = previous

			return nil
		}

		if err := x.complete(ctx, r, event, fingerprint); err != nil {
			return err
		}

		completed = event

		return nil
	}); err != nil {
		return Response{}, err
	}

	return NewResponse(completed), nil
}

func (x service) complete(ctx context.Context, r Request, event *aggregate_event.Event, fingerprint string) error {
	if event.GetCoordinatorID() != r.GetUserID() {
		return ErrPermissionDenied
	}

	now := time.Now().UTC()

	command := aggregate_event_command.NewComplete(r.GetRequestID(), fingerprint, now)
	if r.IsForced() {
		command = aggregate_event_command.NewForceComplete(r.GetRequestID(), fingerprint, now)
	}

	if err := event.Complete(command); err != nil {
		return err
	}

	participations, err := x.participations.List(ctx)
	if err != nil {
		return err
	}

	results := make(map[model_participation_id.ParticipationID]model_event_result.EventResult, len(r.GetResults()))
	for _, result := range r.GetResults() {
		results[result.GetParticipationID()] = result
	}

	approved := make([]*aggregate_participation.Participation, 0)

	for _, value := range participations {
		if value.GetEventID() != event.GetEventID() {
			continue
		}

		if value.GetStatus() == model_participation_status.Approved {
			approved = append(approved, value)
		}
	}

	if len(approved) != len(results) {
		return ErrInvalidRequest
	}

	// Acquire wallet locks in the same order across independent transactions.
	sort.Slice(approved, func(i, j int) bool {
		return approved[i].GetVolunteerID().String() < approved[j].GetVolunteerID().String()
	})

	for _, value := range approved {
		result, ok := results[value.GetParticipationID()]
		if !ok {
			return ErrInvalidRequest
		}

		if err := x.completeParticipation(ctx, event, value, result, now); err != nil {
			return err
		}
	}

	for _, value := range participations {
		if value.GetEventID() != event.GetEventID() || value.GetStatus() != model_participation_status.Pending {
			continue
		}

		previous := utils_value.NewOptional(value.GetStatus())
		if err := value.Close(aggregate_participation_command.NewClose(now)); err != nil {
			return err
		}

		if err := x.participations.Save(ctx, value); err != nil {
			return err
		}

		if err := x.recordChange(ctx, value, previous, model_user_id.UserID{}, model_actor_kind.System); err != nil {
			return err
		}
	}

	return nil
}

func (x service) completeParticipation(
	ctx context.Context,
	event *aggregate_event.Event,
	value *aggregate_participation.Participation,
	result model_event_result.EventResult,
	now time.Time,
) error {
	var (
		err         error
		operationID model_operation_id.OperationID
		coins       int64
	)

	if !result.IsAbsent() {
		operationID, err = model_operation_id.GenerateOperationID()
		if err != nil {
			return err
		}

		if err := x.wallets.Modify(ctx, value.GetVolunteerID(), func(
			ctx context.Context,
			wallet *aggregate_wallet.Wallet,
		) error {
			operation, err := wallet.Accrue(aggregate_wallet_command.NewAccrue(
				operationID,
				event.GetCoordinatorID(),
				result,
				event.GetTitle().String(),
				now,
			))
			if err != nil {
				return err
			}

			coins = operation.GetCoinDelta()

			return x.wallets.AppendOperation(ctx, operation)
		}); err != nil {
			return err
		}
	}

	previous := utils_value.NewOptional(value.GetStatus())
	if err := value.Complete(aggregate_participation_command.NewComplete(result, coins, operationID, now)); err != nil {
		return err
	}

	if err := x.participations.Save(ctx, value); err != nil {
		return err
	}

	if err := x.recordChange(ctx, value, previous, event.GetCoordinatorID(), model_actor_kind.Coordinator); err != nil {
		return err
	}

	return nil
}

func (x service) recordChange(
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
		value.GetCancellationReason(),
		value.GetRejectionReason(),
		0,
		value.GetUpdatedAt(),
	)
	if err != nil {
		return err
	}

	return x.participations.AppendChange(ctx, change)
}

func resultFingerprint(results []model_event_result.EventResult) string {
	parts := make([]string, len(results))
	for index, result := range results {
		parts[index] = fmt.Sprintf("%s:%d", result.GetParticipationID().String(), result.GetMinutes())
	}

	sort.Strings(parts)

	return fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join(parts, "|"))))
}

func NewService(
	events aggregate_event.Repository,
	coordinators aggregate_coordinator.Repository,
	participations aggregate_participation.Repository,

	wallets aggregate_wallet.Repository,
) Service {
	return service{
		events:         events,
		coordinators:   coordinators,
		participations: participations,

		wallets: wallets,
	}
}
