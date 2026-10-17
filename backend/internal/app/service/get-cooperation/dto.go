package service_get_cooperation

import (
	"fmt"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID        model_user_id.UserID
	counterpartID model_user_id.UserID
	asCoordinator bool
}

func (x Request) GetAsCoordinator() bool {
	return x.asCoordinator
}

func (x Request) GetCounterpartID() model_user_id.UserID {
	return x.counterpartID
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if err := x.counterpartID.Validate(); err != nil {
		return fmt.Errorf("counterpart id: %w", err)
	}

	return nil
}

func NewRequest(userID, counterpartID model_user_id.UserID, asCoordinator bool) Request {
	return Request{
		userID:        userID,
		counterpartID: counterpartID,
		asCoordinator: asCoordinator,
	}
}

type Response struct {
	cooperation *aggregate_cooperation.Cooperation
}

func (x Response) GetCooperation() *aggregate_cooperation.Cooperation {
	return x.cooperation
}

func NewResponse(cooperation *aggregate_cooperation.Cooperation) Response {
	return Response{
		cooperation: cooperation,
	}
}
