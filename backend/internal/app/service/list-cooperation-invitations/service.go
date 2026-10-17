package service_list_cooperation_invitations

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	aggregate_cooperation_invitation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	"github.com/Roum1212/hahaton/internal/pkg/page"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	Service interface {
		ListCooperationInvitations(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		invitations  aggregate_cooperation_invitation.Repository
		coordinators aggregate_coordinator.Repository
	}
)

func (x service) ListCooperationInvitations(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	if _, err := x.coordinators.Load(ctx, r.GetUserID()); err != nil {
		switch {
		case errors.Is(err, aggregate_coordinator.ErrNotFound):
			return Response{}, ErrPermissionDenied
		default:
			return Response{}, err
		}
	}

	values, err := x.invitations.List(ctx)
	if err != nil {
		return Response{}, err
	}

	filtered := make([]*aggregate_cooperation_invitation.CooperationInvitation, 0)
	now := time.Now().UTC()

	for _, value := range values {
		if value.GetCoordinatorID() != r.GetUserID() {
			continue
		}

		if !r.GetIncludeExpired() && !value.IsActive(now) {
			continue
		}

		filtered = append(filtered, value)
	}

	slices.SortFunc(filtered, func(
		a *aggregate_cooperation_invitation.CooperationInvitation,
		b *aggregate_cooperation_invitation.CooperationInvitation,
	) int {
		if value := b.GetCreatedAt().Compare(a.GetCreatedAt()); value != 0 {
			return value
		}

		return strings.Compare(b.GetInvitationID().String(), a.GetInvitationID().String())
	})
	values, total := page.Slice(filtered, r.GetPage())

	return NewResponse(values, total), nil
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
