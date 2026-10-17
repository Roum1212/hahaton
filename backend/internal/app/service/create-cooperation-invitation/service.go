package service_create_cooperation_invitation

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	aggregate_cooperation_invitation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation"
	aggregate_cooperation_invitation_command "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation/command"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	model_invitation_id "github.com/Roum1212/hahaton/internal/domain/model/invitation-id"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	Service interface {
		CreateCooperationInvitation(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		invitations  aggregate_cooperation_invitation.Repository
		coordinators aggregate_coordinator.Repository
	}
)

func (x service) CreateCooperationInvitation(ctx context.Context, r Request) (Response, error) {
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

	now := time.Now().UTC()

	expiresAt := time.Time{}
	if value, ok := r.GetExpiresAt().Get(); ok {
		expiresAt = value
	}

	if _, ok := r.GetExpiresAt().Get(); ok && !expiresAt.After(now) {
		return Response{}, ErrInvalidRequest
	}

	id, err := model_invitation_id.GenerateInvitationID()
	if err != nil {
		return Response{}, err
	}

	const tokenBytes = 32

	secret := make([]byte, tokenBytes)
	if _, err := rand.Read(secret); err != nil {
		return Response{}, err
	}

	invitation := aggregate_cooperation_invitation.NewCooperationInvitation(id, r.GetUserID())
	invitation.Create(aggregate_cooperation_invitation_command.NewCreate(
		base64.RawURLEncoding.EncodeToString(secret),
		expiresAt,
		now,
	))

	if err := x.invitations.Create(ctx, invitation); err != nil {
		return Response{}, err
	}

	return NewResponse(invitation), nil
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
