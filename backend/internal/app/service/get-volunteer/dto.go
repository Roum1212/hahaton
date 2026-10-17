package service_get_volunteer

import (
	"fmt"

	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	aggregate_wallet "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

type Request struct {
	userID model_user_id.UserID

	id utils_value.Optional[model_user_id.UserID]
}

func (x Request) GetID() utils_value.Optional[model_user_id.UserID] {
	return x.id
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if value, ok := x.id.Get(); ok {
		if err := value.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type RequestOption func(x *Request)

func WithID(value model_user_id.UserID) RequestOption {
	return func(x *Request) {
		x.id = utils_value.NewOptional(value)
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
	volunteer *aggregate_volunteer.Volunteer
	wallet    *aggregate_wallet.Wallet
}

func (x Response) GetVolunteer() *aggregate_volunteer.Volunteer {
	return x.volunteer
}

func (x Response) GetWallet() *aggregate_wallet.Wallet {
	return x.wallet
}

func NewResponse(volunteer *aggregate_volunteer.Volunteer, wallet *aggregate_wallet.Wallet) Response {
	return Response{
		volunteer: volunteer,
		wallet:    wallet,
	}
}
