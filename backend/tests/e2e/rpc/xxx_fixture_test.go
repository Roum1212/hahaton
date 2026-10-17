//go:build e2e

package rpc

import (
	"crypto/rand"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	admin_v1 "github.com/Roum1212/hahaton/genproto/admin/v1"
	cooperation_v1 "github.com/Roum1212/hahaton/genproto/cooperation/v1"
	event_v1 "github.com/Roum1212/hahaton/genproto/event/v1"
	geo_v1 "github.com/Roum1212/hahaton/genproto/geo/v1"
	participation_v1 "github.com/Roum1212/hahaton/genproto/participation/v1"
	session_v1 "github.com/Roum1212/hahaton/genproto/session/v1"
)

type actor struct {
	id     string
	cookie *http.Cookie
}

func (x *RPCTestSuite) cookie(header http.Header) *http.Cookie {
	cookie, err := http.ParseSetCookie(header.Get("Set-Cookie"))
	x.Require().NoError(err)
	x.Require().True(cookie.HttpOnly)
	x.Require().True(cookie.Secure)
	x.Require().Equal("/", cookie.Path)
	x.Require().Equal(http.SameSiteNoneMode, cookie.SameSite)

	return cookie
}

func (x *RPCTestSuite) cooperate(coordinator, volunteer actor) {
	_, err := x.clients.cooperation.RequestCooperation(x.T().Context(), request(volunteer,
		&cooperation_v1.RequestCooperationRequest{CoordinatorId: coordinator.id}))
	x.Require().NoError(err)
	_, err = x.clients.cooperation.ApproveCooperation(x.T().Context(), request(coordinator,
		&cooperation_v1.ApproveCooperationRequest{VolunteerId: volunteer.id}))
	x.Require().NoError(err)
}

func (x *RPCTestSuite) coordinator() actor {
	a := x.login()
	r := connect.NewRequest(
		&admin_v1.CreateCoordinatorRequest{
			VolunteerId:                       a.id,
			CooperationRequestCooldownSeconds: proto.Int64(60),
		},
	)
	r.Header().Set("Authorization", "Bearer "+adminAPIKey)
	response, err := x.clients.admin.CreateCoordinator(x.T().Context(), r)
	x.Require().NoError(err)
	x.Require().Equal(a.id, response.Msg.GetId())

	return a
}

func (x *RPCTestSuite) event(coordinator actor, limit int32) string {
	start := time.Now().Add(time.Hour).Truncate(time.Second)
	response, err := x.clients.event.CreateEvent(x.T().Context(), request(coordinator, &event_v1.CreateEventRequest{
		Title:          "Городская акция",
		Description:    "Учебное мероприятие",
		Location:       "Обнинск, парк",
		Point:          point(),
		StartsAt:       timestamppb.New(start),
		EndsAt:         timestamppb.New(start.Add(3 * time.Hour)),
		VolunteerLimit: limit,
	}))
	x.Require().NoError(err)

	return response.Msg.GetId()
}

func (x *RPCTestSuite) finishTime(id string) {
	// Only the clock fixture bypasses RPC: completion must not wait three hours.
	start := time.Now().Add(-4 * time.Hour).Truncate(time.Second)
	tag, err := x.pool.Exec(x.T().Context(), "UPDATE events SET starts_at=$1, ends_at=$2 WHERE id=$3",
		start, start.Add(3*time.Hour), id)
	x.Require().NoError(err)
	x.Require().EqualValues(1, tag.RowsAffected())
}

func (x *RPCTestSuite) login() actor {
	subject := strings.ToLower(rand.Text())

	return x.loginSubject(subject, nil)
}

func (x *RPCTestSuite) loginSubject(subject string, overrides map[string]any) actor {
	opened, err := x.clients.session.OpenSession(x.T().Context(), connect.NewRequest(&emptypb.Empty{}))
	x.Require().NoError(err)
	a := actor{cookie: x.cookie(opened.Header())}
	authenticated, err := x.clients.session.AuthenticateSession(x.T().Context(), request(a,
		&session_v1.AuthenticateSessionRequest{IdToken: x.provider.token(x.T(), subject, overrides)}))
	x.Require().NoError(err)
	a.cookie = x.cookie(authenticated.Header())
	session, err := x.clients.session.GetSession(x.T().Context(), request(a, &emptypb.Empty{}))
	x.Require().NoError(err)

	a.id = session.Msg.GetSession().GetUserId()
	x.Require().NotEmpty(a.id)

	return a
}

func (x *RPCTestSuite) participate(coordinator, volunteer actor, eventID string) string {
	response, err := x.clients.participation.RequestParticipation(x.T().Context(), request(volunteer,
		&participation_v1.RequestParticipationRequest{EventId: eventID}))
	x.Require().NoError(err)

	id := response.Msg.GetId()
	_, err = x.clients.participation.ApproveParticipation(x.T().Context(), request(coordinator,
		&participation_v1.ApproveParticipationRequest{Id: id}))
	x.Require().NoError(err)

	return id
}

func callError(t require.TestingT, err error, code connect.Code) {
	require.Error(t, err)
	require.Equal(t, code, connect.CodeOf(err), "%v", err)
}

func point() *geo_v1.GeoPoint {
	return &geo_v1.GeoPoint{
		Latitude:  proto.Float64(55.1),
		Longitude: proto.Float64(36.6),
	}
}

func request[T any](a actor, message *T) *connect.Request[T] {
	r := connect.NewRequest(message)
	if a.cookie != nil {
		r.Header().Set("Cookie", a.cookie.String())
	}

	return r
}
