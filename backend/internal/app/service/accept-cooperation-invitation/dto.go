package service_accept_cooperation_invitation

import (
	"errors"
	"fmt"
	"strings"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID model_user_id.UserID
	token  string
}

func (x Request) GetToken() string {
	return x.token
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if len(strings.TrimSpace(x.token)) == 0 {
		return errors.New("invitation token is required")
	}

	return nil
}

func NewRequest(userID model_user_id.UserID, token string) Request {
	return Request{
		userID: userID,
		token:  token,
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
