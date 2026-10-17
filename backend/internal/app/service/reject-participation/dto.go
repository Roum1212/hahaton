package service_reject_participation

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

type Request struct {
	userID          model_user_id.UserID
	participationID model_participation_id.ParticipationID

	reason utils_value.Optional[string]
}

func (x Request) GetParticipationID() model_participation_id.ParticipationID {
	return x.participationID
}

func (x Request) GetReason() utils_value.Optional[string] {
	return x.reason
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if err := x.participationID.Validate(); err != nil {
		return fmt.Errorf("participation id: %w", err)
	}

	if value, ok := x.reason.Get(); ok && (len(strings.TrimSpace(value)) == 0 ||
		utf8.RuneCountInString(strings.TrimSpace(value)) > 2000) {
		return errors.New("rejection reason must contain 1 to 2000 characters")
	}

	return nil
}

type RequestOption func(x *Request)

func WithReason(value string) RequestOption {
	return func(x *Request) {
		x.reason = utils_value.NewOptional(value)
	}
}

func NewRequest(
	userID model_user_id.UserID,
	participationID model_participation_id.ParticipationID,
	opts ...RequestOption,
) Request {
	x := Request{
		userID:          userID,
		participationID: participationID,
	}
	for _, option := range opts {
		option(&x)
	}

	return x
}
