//go:build e2e

package rpc

import (
	"crypto/rand"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	session_v1 "github.com/Roum1212/hahaton/genproto/session/v1"
	volunteer_v1 "github.com/Roum1212/hahaton/genproto/volunteer/v1"
)

func (x *RPCTestSuite) TestSessionLifecycle() {
	ctx := x.T().Context()
	_, err := x.clients.session.GetSession(ctx, request(actor{}, &emptypb.Empty{}))
	callError(x.T(), err, connect.CodeUnauthenticated)
	opened, err := x.clients.session.OpenSession(ctx, request(actor{}, &emptypb.Empty{}))
	x.Require().NoError(err)
	anonymous := actor{cookie: x.cookie(opened.Header())}
	initial, err := x.clients.session.GetSession(ctx, request(anonymous, &emptypb.Empty{}))
	x.Require().NoError(err)
	x.Require().Nil(initial.Msg.GetSession().UserId)
	_, err = x.clients.volunteer.GetVolunteer(ctx, request(anonymous, &volunteer_v1.GetVolunteerRequest{}))
	callError(x.T(), err, connect.CodeUnauthenticated)

	subject := strings.ToLower(rand.Text())
	token := x.provider.token(x.T(), subject, nil)
	authenticated, err := x.clients.session.AuthenticateSession(ctx, request(anonymous,
		&session_v1.AuthenticateSessionRequest{IdToken: token}))
	x.Require().NoError(err)
	a := actor{cookie: x.cookie(authenticated.Header())}
	x.Require().NotEqual(anonymous.cookie.Value, a.cookie.Value)
	_, err = x.clients.session.GetSession(ctx, request(anonymous, &emptypb.Empty{}))
	callError(x.T(), err, connect.CodeUnauthenticated)
	session, err := x.clients.session.GetSession(ctx, request(a, &emptypb.Empty{}))
	x.Require().NoError(err)

	a.id = session.Msg.GetSession().GetUserId()
	x.Require().NotEmpty(a.id)
	x.Require().False(session.Msg.GetSession().GetHasCoordinatorProfile())
	reopened, err := x.clients.session.OpenSession(ctx, request(a, &emptypb.Empty{}))
	x.Require().NoError(err)
	retained, err := x.clients.session.GetSession(
		ctx,
		request(actor{cookie: x.cookie(reopened.Header())}, &emptypb.Empty{}),
	)
	x.Require().NoError(err)
	x.Require().Equal(a.id, retained.Msg.GetSession().GetUserId())
	_, err = x.clients.session.AuthenticateSession(ctx, request(a,
		&session_v1.AuthenticateSessionRequest{IdToken: token}))
	callError(x.T(), err, connect.CodeUnauthenticated)
	other := x.loginSubject(subject, nil)
	x.Require().Equal(a.id, other.id)
	closed, err := x.clients.session.CloseSession(ctx, request(a, &emptypb.Empty{}))
	x.Require().NoError(err)
	cookie, err := http.ParseSetCookie(closed.Header().Get("Set-Cookie"))
	x.Require().NoError(err)
	x.Require().Negative(cookie.MaxAge)
	_, err = x.clients.session.GetSession(ctx, request(a, &emptypb.Empty{}))
	callError(x.T(), err, connect.CodeUnauthenticated)
	_, err = x.clients.session.GetSession(ctx, request(other, &emptypb.Empty{}))
	x.Require().NoError(err)
	_, err = x.clients.session.CloseSession(ctx, request(actor{}, &emptypb.Empty{}))
	x.Require().NoError(err)
}

func (x *RPCTestSuite) TestSessionRejectsInvalidTokens() {
	opened, err := x.clients.session.OpenSession(x.T().Context(), request(actor{}, &emptypb.Empty{}))
	x.Require().NoError(err)
	anonymous := actor{cookie: x.cookie(opened.Header())}

	cases := []struct {
		name   string
		claims map[string]any
	}{
		{
			name:   "audience",
			claims: map[string]any{"aud": "other"},
		},
		{
			name:   "issuer",
			claims: map[string]any{"iss": "https://other.test"},
		},
		{
			name:   "expired",
			claims: map[string]any{"exp": time.Now().Add(-time.Hour).Unix()},
		},
		{
			name:   "future issuance",
			claims: map[string]any{"iat": time.Now().Add(time.Hour).Unix()},
		},
		{
			name:   "not before",
			claims: map[string]any{"nbf": time.Now().Add(time.Hour).Unix()},
		},
		{
			name:   "subject",
			claims: map[string]any{"sub": ""},
		},
		{
			name:   "authorized party",
			claims: map[string]any{"azp": "other"},
		},
	}
	for _, scenario := range cases {
		x.Run(scenario.name, func() {
			_, err := x.clients.session.AuthenticateSession(x.T().Context(), request(anonymous,
				&session_v1.AuthenticateSessionRequest{IdToken: x.provider.token(x.T(), "invalid", scenario.claims)}))
			callError(x.T(), err, connect.CodeUnauthenticated)
		})
	}

	_, err = x.clients.session.AuthenticateSession(x.T().Context(), request(anonymous,
		&session_v1.AuthenticateSessionRequest{IdToken: "invalid"}))
	callError(x.T(), err, connect.CodeUnauthenticated)
	_, err = x.clients.session.AuthenticateSession(x.T().Context(), request(anonymous,
		&session_v1.AuthenticateSessionRequest{}))
	callError(x.T(), err, connect.CodeInvalidArgument)
}
