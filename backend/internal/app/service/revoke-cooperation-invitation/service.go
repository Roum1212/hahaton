package service_revoke_cooperation_invitation

import (
	"context"
	"errors"
	"fmt"

	aggregate_cooperation_invitation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	Service interface {
		RevokeCooperationInvitation(ctx context.Context, r Request) error
	}
	service struct {
		invitations  aggregate_cooperation_invitation.Repository
		coordinators aggregate_coordinator.Repository
	}
)

func (x service) RevokeCooperationInvitation(ctx context.Context, r Request) error {
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

	invitation, err := x.invitations.Load(ctx, r.GetInvitationID())
	if err != nil {
		switch {
		case errors.Is(err, aggregate_cooperation_invitation.ErrNotFound):
			return nil
		default:
			return err
		}
	}

	if invitation.GetCoordinatorID() != r.GetUserID() {
		return ErrPermissionDenied
	}

	return x.invitations.Delete(ctx, r.GetInvitationID())
}

func NewService(
	invitations aggregate_cooperation_invitation.Repository,
	coordinators aggregate_coordinator.Repository,
) Service {
	return service{
		invitations:  invitations,
		coordinators: coordinators,
	}
}
