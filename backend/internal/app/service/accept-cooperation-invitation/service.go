package service_accept_cooperation_invitation

import (
	"context"
	"errors"
	"fmt"
	"time"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	aggregate_cooperation_invitation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation"
	aggregate_cooperation_command "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation/command"
	model_cooperation_change "github.com/Roum1212/hahaton/internal/domain/model/cooperation-change"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

var (
	ErrFailedPrecondition = errors.New("failed precondition")
	ErrInvalidRequest     = errors.New("invalid request")
)

type (
	Service interface {
		AcceptCooperationInvitation(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		invitations  aggregate_cooperation_invitation.Repository
		cooperations aggregate_cooperation.Repository
	}
)

func (x service) AcceptCooperationInvitation(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	invitations, err := x.invitations.List(ctx)
	if err != nil {
		return Response{}, err
	}

	var invitation *aggregate_cooperation_invitation.CooperationInvitation

	for _, value := range invitations {
		if value.GetToken() == r.GetToken() {
			invitation = value

			break
		}
	}

	if invitation == nil {
		return Response{}, aggregate_cooperation_invitation.ErrNotFound
	}

	var result *aggregate_cooperation.Cooperation

	if err := x.invitations.Modify(ctx, invitation.GetInvitationID(), func(
		ctx context.Context,
		invitation *aggregate_cooperation_invitation.CooperationInvitation,
	) error {
		if !invitation.IsActive(time.Now().UTC()) {
			return ErrFailedPrecondition
		}

		value, err := x.cooperations.Load(ctx, invitation.GetCoordinatorID(), r.GetUserID())
		exists := true

		if err != nil {
			switch {
			case errors.Is(err, aggregate_cooperation.ErrNotFound):
				value = aggregate_cooperation.NewCooperation(invitation.GetCoordinatorID(), r.GetUserID())
				exists = false
			default:
				return err
			}
		}

		previousStatus := value.GetStatus()
		if err := value.AcceptInvitation(
			aggregate_cooperation_command.NewAcceptInvitation(time.Now().UTC()),
		); err != nil {
			return err
		}

		if exists {
			err = x.cooperations.Save(ctx, value)
		} else {
			err = x.cooperations.Create(ctx, value, func(ctx context.Context) error {
				return nil
			})
		}

		if err != nil {
			return err
		}

		if value.GetStatus() != previousStatus {
			if err := x.recordCooperationChange(ctx, value, r.GetUserID()); err != nil {
				return err
			}
		}

		result = value

		return nil
	}); err != nil {
		return Response{}, err
	}

	return NewResponse(result), nil
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

func NewService(
	invitations aggregate_cooperation_invitation.Repository,
	cooperations aggregate_cooperation.Repository,
) Service {
	return service{
		invitations:  invitations,
		cooperations: cooperations,
	}
}
