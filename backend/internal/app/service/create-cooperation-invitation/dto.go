package service_create_cooperation_invitation

import (
	"fmt"
	"time"

	aggregate_cooperation_invitation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

type Request struct {
	userID model_user_id.UserID

	expiresAt utils_value.Optional[time.Time]
}

func (x Request) GetExpiresAt() utils_value.Optional[time.Time] {
	return x.expiresAt
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	return nil
}

type RequestOption func(x *Request)

func WithExpiresAt(value time.Time) RequestOption {
	return func(x *Request) {
		x.expiresAt = utils_value.NewOptional(value)
	}
}

func NewRequest(userID model_user_id.UserID, opts ...RequestOption) Request {
	x := Request{
		userID: userID,
	}
	for _, option := range opts {
		option(&x)
	}

	return x
}

type Response struct {
	invitation *aggregate_cooperation_invitation.CooperationInvitation
}

func (x Response) GetInvitation() *aggregate_cooperation_invitation.CooperationInvitation {
	return x.invitation
}

func NewResponse(invitation *aggregate_cooperation_invitation.CooperationInvitation) Response {
	return Response{
		invitation: invitation,
	}
}
