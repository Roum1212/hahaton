package service_get_coordinator

import (
	"fmt"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
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
	coordinator *aggregate_coordinator.Coordinator
	volunteer   *aggregate_volunteer.Volunteer
}

func (x Response) GetCoordinator() *aggregate_coordinator.Coordinator {
	return x.coordinator
}

func (x Response) GetVolunteer() *aggregate_volunteer.Volunteer {
	return x.volunteer
}

func NewResponse(coordinator *aggregate_coordinator.Coordinator, volunteer *aggregate_volunteer.Volunteer) Response {
	return Response{
		coordinator: coordinator,
		volunteer:   volunteer,
	}
}
