package service_update_coordinator_settings

import (
	"errors"
	"fmt"

	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

type Request struct {
	userID model_user_id.UserID

	cooldownSeconds utils_value.Optional[int64]
}

func (x Request) GetCooldownSeconds() utils_value.Optional[int64] {
	return x.cooldownSeconds
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if value, ok := x.cooldownSeconds.Get(); ok && value < 0 {
		return errors.New("cooldown must not be negative")
	}

	return nil
}

type RequestOption func(x *Request)

func WithCooldownSeconds(value int64) RequestOption {
	return func(x *Request) {
		x.cooldownSeconds = utils_value.NewOptional(value)
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
