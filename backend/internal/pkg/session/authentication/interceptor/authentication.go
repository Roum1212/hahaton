package session_authentication_interceptor

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/authn"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/Roum1212/hahaton/genproto/cooperation/v1/cooperation_v1connect"
	"github.com/Roum1212/hahaton/genproto/coordinator/v1/coordinator_v1connect"
	"github.com/Roum1212/hahaton/genproto/event-chat/v1/event_chat_v1connect"
	"github.com/Roum1212/hahaton/genproto/event/v1/event_v1connect"
	"github.com/Roum1212/hahaton/genproto/geo/v1/geo_v1connect"
	"github.com/Roum1212/hahaton/genproto/image/v1/image_v1connect"
	"github.com/Roum1212/hahaton/genproto/participation/v1/participation_v1connect"
	"github.com/Roum1212/hahaton/genproto/reward-order/v1/reward_order_v1connect"
	"github.com/Roum1212/hahaton/genproto/reward/v1/reward_v1connect"
	"github.com/Roum1212/hahaton/genproto/session/v1/session_v1connect"
	"github.com/Roum1212/hahaton/genproto/volunteer/v1/volunteer_v1connect"
	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
	"github.com/Roum1212/hahaton/internal/pkg/cookies"
	rpc_error "github.com/Roum1212/hahaton/internal/pkg/rpc-error"
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
	if strings.HasPrefix(r.URL.Path, "/storage/") || r.URL.Path == "/healthz" {
		return nil, nil
	}

	procedure, _ := authn.InferProcedure(r.URL)

	mode, ok := authModeByProcedure(procedure)
	if !ok {
		return nil, rpc_error.NewNotFoundError(&errdetails.ResourceInfo{
			ResourceType: "connect.Procedure",
			ResourceName: r.URL.Path,
			Description:  "procedure not found",
		})
	}

	if mode == authModePublic {
		return nil, nil
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
	case authModePublic:
		return nil, nil
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
	authModePublic
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
	modes := map[string]authMode{
		session_v1connect.SessionServiceAuthenticateSessionProcedure: authModeUnauthenticated,
		session_v1connect.SessionServiceCloseSessionProcedure:        authModeActiveSession,
		session_v1connect.SessionServiceGetSessionProcedure:          authModeActiveSession,
		session_v1connect.SessionServiceOpenSessionProcedure:         authModeActiveSession,

		cooperation_v1connect.CooperationServiceAcceptCooperationInvitationProcedure: authModeAuthenticated,
		cooperation_v1connect.CooperationServiceApproveCooperationProcedure:          authModeAuthenticated,
		cooperation_v1connect.CooperationServiceCreateCooperationInvitationProcedure: authModeAuthenticated,
		cooperation_v1connect.CooperationServiceGetCooperationProcedure:              authModeAuthenticated,
		cooperation_v1connect.CooperationServiceListCooperationInvitationsProcedure:  authModeAuthenticated,
		cooperation_v1connect.CooperationServiceListCooperationsProcedure:            authModeAuthenticated,
		cooperation_v1connect.CooperationServiceRejectCooperationProcedure:           authModeAuthenticated,
		cooperation_v1connect.CooperationServiceRequestCooperationProcedure:          authModeAuthenticated,
		cooperation_v1connect.CooperationServiceRevokeCooperationInvitationProcedure: authModeAuthenticated,
		cooperation_v1connect.CooperationServiceTerminateCooperationProcedure:        authModeAuthenticated,

		coordinator_v1connect.CoordinatorServiceGetCoordinatorProcedure:            authModeAuthenticated,
		coordinator_v1connect.CoordinatorServiceGetCoordinatorSettingsProcedure:    authModeAuthenticated,
		coordinator_v1connect.CoordinatorServiceListCoordinatorsProcedure:          authModeAuthenticated,
		coordinator_v1connect.CoordinatorServiceUpdateCoordinatorProcedure:         authModeAuthenticated,
		coordinator_v1connect.CoordinatorServiceUpdateCoordinatorSettingsProcedure: authModeAuthenticated,

		event_chat_v1connect.EventChatServiceGetEventChatProcedure:               authModeAuthenticated,
		event_chat_v1connect.EventChatServiceListEventChatMessagesProcedure:      authModeAuthenticated,
		event_chat_v1connect.EventChatServiceListMyEventChatsProcedure:           authModeAuthenticated,
		event_chat_v1connect.EventChatServiceSendEventChatMessageProcedure:       authModeAuthenticated,
		event_chat_v1connect.EventChatServiceSubscribeEventChatMessagesProcedure: authModeAuthenticated,

		event_v1connect.EventServiceCancelEventProcedure:   authModeAuthenticated,
		event_v1connect.EventServiceCompleteEventProcedure: authModeAuthenticated,
		event_v1connect.EventServiceCreateEventProcedure:   authModeAuthenticated,
		event_v1connect.EventServiceGetEventProcedure:      authModeAuthenticated,
		event_v1connect.EventServiceListEventsProcedure:    authModeAuthenticated,
		event_v1connect.EventServiceUpdateEventProcedure:   authModeAuthenticated,

		geo_v1connect.GeoServiceGetCityProcedure:          authModePublic,
		geo_v1connect.GeoServiceListCitiesProcedure:       authModePublic,
		geo_v1connect.GeoServiceSuggestAddressesProcedure: authModeAuthenticated,

		image_v1connect.ImageServiceUploadImageProcedure: authModeAuthenticated,

		participation_v1connect.ParticipationServiceApproveParticipationProcedure:     authModeAuthenticated,
		participation_v1connect.ParticipationServiceCancelParticipationProcedure:      authModeAuthenticated,
		participation_v1connect.ParticipationServiceGetParticipationProcedure:         authModeAuthenticated,
		participation_v1connect.ParticipationServiceListParticipationChangesProcedure: authModeAuthenticated,
		participation_v1connect.ParticipationServiceListParticipationsProcedure:       authModeAuthenticated,
		participation_v1connect.ParticipationServiceRejectParticipationProcedure:      authModeAuthenticated,
		participation_v1connect.ParticipationServiceRequestParticipationProcedure:     authModeAuthenticated,
		participation_v1connect.ParticipationServiceRestoreParticipationProcedure:     authModeAuthenticated,

		reward_order_v1connect.RewardOrderServiceCancelRewardOrderProcedure: authModeAuthenticated,
		reward_order_v1connect.RewardOrderServiceCreateRewardOrderProcedure: authModeAuthenticated,
		reward_order_v1connect.RewardOrderServiceGetRewardOrderProcedure:    authModeAuthenticated,
		reward_order_v1connect.RewardOrderServiceIssueRewardOrderProcedure:  authModeAuthenticated,
		reward_order_v1connect.RewardOrderServiceListRewardOrdersProcedure:  authModeAuthenticated,

		reward_v1connect.RewardServiceCreateRewardProcedure: authModeAuthenticated,
		reward_v1connect.RewardServiceDeleteRewardProcedure: authModeAuthenticated,
		reward_v1connect.RewardServiceGetRewardProcedure:    authModeAuthenticated,
		reward_v1connect.RewardServiceListRewardsProcedure:  authModeAuthenticated,
		reward_v1connect.RewardServiceUpdateRewardProcedure: authModeAuthenticated,

		volunteer_v1connect.VolunteerServiceGetCoordinatorVolunteerProcedure: authModeAuthenticated,
		volunteer_v1connect.VolunteerServiceGetVolunteerProcedure:            authModeAuthenticated,
		volunteer_v1connect.VolunteerServiceListVolunteerOperationsProcedure: authModeAuthenticated,
		volunteer_v1connect.VolunteerServiceUpdateVolunteerProcedure:         authModeAuthenticated,
	}
	mode, ok := modes[procedure]

	return mode, ok
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
