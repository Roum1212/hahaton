package service_revoke_cooperation_invitation

import (
	"fmt"

	model_invitation_id "github.com/Roum1212/hahaton/internal/domain/model/invitation-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID       model_user_id.UserID
	invitationID model_invitation_id.InvitationID
}

func (x Request) GetInvitationID() model_invitation_id.InvitationID {
	return x.invitationID
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if err := x.invitationID.Validate(); err != nil {
		return fmt.Errorf("invitation id: %w", err)
	}

	return nil
}

func NewRequest(userID model_user_id.UserID, invitationID model_invitation_id.InvitationID) Request {
	return Request{
		userID:       userID,
		invitationID: invitationID,
	}
}
