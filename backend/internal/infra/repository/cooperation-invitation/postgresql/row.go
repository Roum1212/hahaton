package repository_cooperation_invitation_postgresql

import (
	"encoding/base64"
	"fmt"
	"time"

	aggregate_cooperation_invitation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation"
	model_invitation_id "github.com/Roum1212/hahaton/internal/domain/model/invitation-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type row struct {
	invitationID  string
	coordinatorID string
	token         string
	createdAt     time.Time
	expiresAt     time.Time
}

func (x *row) scan(source interface{ Scan(dest ...any) error }) error {
	return source.Scan(&x.invitationID, &x.coordinatorID, &x.token, &x.createdAt, &x.expiresAt)
}

func (x *row) toAggregate(crypto Crypto) (*aggregate_cooperation_invitation.CooperationInvitation, error) {
	invitationID, err := model_invitation_id.NewInvitationID(x.invitationID)
	if err != nil {
		return nil, fmt.Errorf("restore invitation id: %w", err)
	}

	coordinatorID, err := model_user_id.NewUserID(x.coordinatorID)
	if err != nil {
		return nil, fmt.Errorf("restore coordinator id: %w", err)
	}

	encrypted, err := base64.RawStdEncoding.DecodeString(x.token)
	if err != nil {
		return nil, err
	}

	token, err := crypto.Decrypt(encrypted)
	if err != nil {
		return nil, err
	}

	return aggregate_cooperation_invitation.NewCooperationInvitationFromFields(
		invitationID,
		coordinatorID,
		token,
		x.createdAt,
		x.expiresAt,
	), nil
}
