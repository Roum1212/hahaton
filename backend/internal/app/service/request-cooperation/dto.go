package service_request_cooperation

import (
	"errors"
	"fmt"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID        model_user_id.UserID
	coordinatorID model_user_id.UserID
}

func (x Request) GetCoordinatorID() model_user_id.UserID {
	return x.coordinatorID
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if err := x.coordinatorID.Validate(); err != nil {
		return fmt.Errorf("coordinator id: %w", err)
	}

	if x.userID == x.coordinatorID {
		return errors.New("cannot cooperate with yourself")
	}

	return nil
}

func NewRequest(userID, coordinatorID model_user_id.UserID) Request {
	return Request{
		userID:        userID,
		coordinatorID: coordinatorID,
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
