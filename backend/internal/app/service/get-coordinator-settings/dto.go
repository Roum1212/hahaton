package service_get_coordinator_settings

import (
	"fmt"

	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID model_user_id.UserID
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

func NewRequest(userID model_user_id.UserID) Request {
	return Request{
		userID: userID,
	}
}

type Response struct {
	cooldownSeconds int64
}

func (x Response) GetCooldownSeconds() int64 {
	return x.cooldownSeconds
}

func NewResponse(cooldownSeconds int64) Response {
	return Response{
		cooldownSeconds: cooldownSeconds,
	}
}
