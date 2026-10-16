package session_authentication_interceptor

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/authn"
	"go.flida.dev/cookies/pkg/cookies"
	rpc_error "go.flida.dev/rpc-error"

	"github.com/Roum1212/hahaton/genproto/image/v1/image_v1connect"
	"github.com/Roum1212/hahaton/genproto/session/v1/session_v1connect"
	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
	session_token "github.com/Roum1212/hahaton/internal/pkg/session/token"
	"github.com/Roum1212/hahaton/internal/pkg/token"
)

type Authentication struct {
	cookieService         cookies.Service
	sessionRepository     aggregate_session.Repository
	sessionTokenDecryptor token.Decryptor[session_token.Token]
}

//nolint:nilnil // authn uses nil authentication info for requests without an applicable session.
func (x Authentication) Authenticate(ctx context.Context, r *http.Request) (any, error) {
	if strings.HasPrefix(r.URL.Path, "/storage/") {
		return nil, nil
	}

	procedure, _ := authn.InferProcedure(r.URL)

	mode, ok := authModeByProcedure(procedure)
	if !ok {
		return nil, rpc_error.NewNotFoundError(nil)
	}

	session, err := x.getActiveSession(ctx, r)
	if err != nil {
		switch {
		case errors.Is(err, aggregate_session.ErrInvalidCredentials):
			return nil, nil
		default:
			return nil, rpc_error.NewInternalError()
		}
	}

	switch mode {
	case authModeUnauthenticated:
		if !session.IsAnonymous() {
			return nil, nil
		}
	case authModeAuthenticated:
		if session.IsAnonymous() {
			return nil, nil
		}
	case authModeActiveSession:
	default:
		return nil, rpc_error.NewInternalError()
	}

	return session, nil
}

type authMode int8

const (
	authModeUnauthenticated authMode = iota
	authModeAuthenticated
	authModeActiveSession
)

func (x Authentication) getActiveSession(ctx context.Context, r *http.Request) (*aggregate_session.Session, error) {
	sessionToken, ok := x.getSessionToken(r)
	if !ok {
		return nil, aggregate_session.ErrInvalidCredentials
	}

	session, err := x.sessionRepository.Load(ctx, sessionToken.GetSessionID())
	if err != nil {
		switch {
		case errors.Is(err, aggregate_session.ErrNotFound):
			return nil, aggregate_session.ErrInvalidCredentials
		default:
			return nil, fmt.Errorf("load session: %w", err)
		}
	}

	if !session.IsActive(time.Now().UTC()) {
		return nil, aggregate_session.ErrInvalidCredentials
	}

	return session, nil
}

func (x Authentication) getSessionToken(r *http.Request) (session_token.Token, bool) {
	cookie, err := x.cookieService.GetCookie(r)
	if err != nil {
		return session_token.Token{}, false
	}

	encryptedToken, err := token.NewEncryptedToken(cookie.Value)
	if err != nil {
		return session_token.Token{}, false
	}

	sessionToken, err := x.sessionTokenDecryptor.Decrypt(encryptedToken)
	if err != nil {
		return session_token.Token{}, false
	}

	return sessionToken, true
}

func authModeByProcedure(procedure string) (authMode, bool) {
	switch procedure {
	case session_v1connect.SessionServiceAuthenticateSessionProcedure:
		return authModeUnauthenticated, true
	case image_v1connect.ImageServiceUploadImageProcedure:
		return authModeAuthenticated, true
	case session_v1connect.SessionServiceCloseSessionProcedure,
		session_v1connect.SessionServiceGetSessionProcedure,
		session_v1connect.SessionServiceOpenSessionProcedure:
		return authModeActiveSession, true
	default:
		return authModeActiveSession, false
	}
}

func NewAuthentication(
	cookieService cookies.Service,
	sessionRepository aggregate_session.Repository,
	sessionTokenDecryptor token.Decryptor[session_token.Token],
) Authentication {
	return Authentication{
		cookieService:         cookieService,
		sessionRepository:     sessionRepository,
		sessionTokenDecryptor: sessionTokenDecryptor,
	}
}
