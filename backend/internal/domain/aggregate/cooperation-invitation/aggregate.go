package aggregate_cooperation_invitation

import (
	"time"

	aggregate_cooperation_invitation_command "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation/command"
	model_invitation_id "github.com/Roum1212/hahaton/internal/domain/model/invitation-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type CooperationInvitation struct {
	invitationID  model_invitation_id.InvitationID
	coordinatorID model_user_id.UserID
	token         string
	createdAt     time.Time
	expiresAt     time.Time
}

func (x *CooperationInvitation) GetCoordinatorID() model_user_id.UserID {
	return x.coordinatorID
}

func (x *CooperationInvitation) GetCreatedAt() time.Time {
	return x.createdAt
}

func (x *CooperationInvitation) GetExpiresAt() time.Time {
	return x.expiresAt
}

func (x *CooperationInvitation) GetInvitationID() model_invitation_id.InvitationID {
	return x.invitationID
}

func (x *CooperationInvitation) GetToken() string {
	return x.token
}

func (x *CooperationInvitation) IsActive(now time.Time) bool {
	return x.expiresAt.IsZero() || now.Before(x.expiresAt)
}

func (x *CooperationInvitation) Create(command aggregate_cooperation_invitation_command.Create) {
	x.token = command.GetToken()
	x.expiresAt = command.GetExpiresAt()
	x.createdAt = command.GetNow()
}

func NewCooperationInvitation(
	invitationID model_invitation_id.InvitationID,
	coordinatorID model_user_id.UserID,
) *CooperationInvitation {
	return &CooperationInvitation{
		invitationID:  invitationID,
		coordinatorID: coordinatorID,
	}
}

func NewCooperationInvitationFromFields(
	invitationID model_invitation_id.InvitationID,
	coordinatorID model_user_id.UserID,
	token string,
	createdAt time.Time,
	expiresAt time.Time,
) *CooperationInvitation {
	return &CooperationInvitation{
		invitationID:  invitationID,
		coordinatorID: coordinatorID,
		token:         token,
		createdAt:     createdAt,
		expiresAt:     expiresAt,
	}
}
