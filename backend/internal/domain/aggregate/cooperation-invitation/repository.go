package aggregate_cooperation_invitation

import (
	"context"
	"errors"

	model_invitation_id "github.com/Roum1212/hahaton/internal/domain/model/invitation-id"
)

var (
	ErrNotFound      = errors.New("cooperation invitation not found")
	ErrAlreadyExists = errors.New("cooperation invitation already exists")
)

type Repository interface {
	Create(ctx context.Context, value *CooperationInvitation) error
	Delete(ctx context.Context, invitationID model_invitation_id.InvitationID) error
	List(ctx context.Context) ([]*CooperationInvitation, error)
	Load(ctx context.Context, invitationID model_invitation_id.InvitationID) (*CooperationInvitation, error)
	Modify(
		ctx context.Context,
		invitationID model_invitation_id.InvitationID,
		change func(
			ctx context.Context,
			value *CooperationInvitation,
		) error,
	) error
	Save(ctx context.Context, value *CooperationInvitation) error
}
