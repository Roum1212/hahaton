package service_get_session

import (
	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
)

type Request struct{}

func (x Request) Validate() error {
	return nil
}

func NewRequest() Request {
	return Request{}
}

type Response struct {
	session               *aggregate_session.Session
	hasCoordinatorProfile bool
}

func (x Response) GetSession() *aggregate_session.Session {
	return x.session
}

func (x Response) HasCoordinatorProfile() bool {
	return x.hasCoordinatorProfile
}

func NewResponse(session *aggregate_session.Session, hasCoordinatorProfile bool) Response {
	return Response{
		session:               session,
		hasCoordinatorProfile: hasCoordinatorProfile,
	}
}
