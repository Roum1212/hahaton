package service_request_cooperation

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

var ErrInvalidRequest = errors.New("invalid request")

type (
	Service interface {
		RequestCooperation(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		cooperations aggregate_cooperation.Repository
		coordinators aggregate_coordinator.Repository
	}
)

func (x service) RequestCooperation(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	if _, err := x.coordinators.Load(ctx, r.GetCoordinatorID()); err != nil {
		return Response{}, err
	}

	value, err := x.cooperations.Load(ctx, r.GetCoordinatorID(), r.GetUserID())
	exists := true

	if err != nil {
		switch {
		case errors.Is(err, aggregate_cooperation.ErrNotFound):
			value = aggregate_cooperation.NewCooperation(r.GetCoordinatorID(), r.GetUserID())
			exists = false
		default:
			return Response{}, err
		}
	}

	change := func(
		ctx context.Context,
		current *aggregate_cooperation.Cooperation,
	) error {
		if err := current.Request(aggregate_cooperation_command.NewRequest(time.Now().UTC())); err != nil {
			return err
		}

		value = current

		return x.recordCooperationChange(ctx, current, r.GetUserID())
	}
	if exists {
		err = x.cooperations.Modify(ctx, r.GetCoordinatorID(), r.GetUserID(), change)
	} else {
		if err := value.Request(aggregate_cooperation_command.NewRequest(time.Now().UTC())); err != nil {
			return Response{}, err
		}

		err = x.cooperations.Create(ctx, value, func(ctx context.Context) error {
			return x.recordCooperationChange(ctx, value, r.GetUserID())
		})
	}

	if err != nil {
		return Response{}, err
	}

	return NewResponse(value), nil
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
