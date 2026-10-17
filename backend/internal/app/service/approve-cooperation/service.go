package service_approve_cooperation

import (
	"context"
	"errors"
	"fmt"
	"time"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_cooperation_command "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation/command"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	model_cooperation_change "github.com/Roum1212/hahaton/internal/domain/model/cooperation-change"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	Service interface {
		ApproveCooperation(ctx context.Context, r Request) error
	}
	service struct {
		cooperations aggregate_cooperation.Repository
		coordinators aggregate_coordinator.Repository
	}
)

func (x service) ApproveCooperation(ctx context.Context, r Request) error {
	if err := r.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	if _, err := x.coordinators.Load(ctx, r.GetUserID()); err != nil {
		switch {
		case errors.Is(err, aggregate_coordinator.ErrNotFound):
			return ErrPermissionDenied
		default:
			return err
		}
	}

	return x.cooperations.Modify(ctx, r.GetUserID(), r.GetVolunteerID(), func(
		ctx context.Context,
		value *aggregate_cooperation.Cooperation,
	) error {
		if err := value.Approve(aggregate_cooperation_command.NewApprove(time.Now().UTC())); err != nil {
			return err
		}

		return x.recordCooperationChange(ctx, value, r.GetUserID())
	})
}

func (x service) recordCooperationChange(
	ctx context.Context,
	value *aggregate_cooperation.Cooperation,
	actorID model_user_id.UserID,
) error {
	change, err := model_cooperation_change.NewCooperationChange(value.GetCoordinatorID(), value.GetVolunteerID(),
		value.GetStatus(), actorID, value.GetRetryAfter(), value.GetUpdatedAt())
	if err != nil {
		return err
	}

	return x.cooperations.AppendChange(ctx, change)
}

func NewService(cooperations aggregate_cooperation.Repository, coordinators aggregate_coordinator.Repository) Service {
	return service{
		cooperations: cooperations,
		coordinators: coordinators,
	}
}
