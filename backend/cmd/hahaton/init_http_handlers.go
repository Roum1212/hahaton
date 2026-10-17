package hahaton

import (
	"net/http"

	"connectrpc.com/authn"
	"connectrpc.com/connect"

	"github.com/Roum1212/hahaton/genproto/admin/v1/admin_v1connect"
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
	"github.com/Roum1212/hahaton/internal/config"
	session_authentication_interceptor "github.com/Roum1212/hahaton/internal/pkg/session/authentication/interceptor"
)

const imageReadLimit = 8 << 20

func initHTTPHandlers(
	cfg config.Config,
	repos repositorySet,
	tokens tokenSet,
	supportServices supportServiceSet,
	rpcServices rpcServiceSet,
) (http.Handler, http.Handler) {
	mux := http.NewServeMux()
	adminMux := http.NewServeMux()
	adminMux.Handle(admin_v1connect.NewAdminServiceHandler(rpcServices.admin))
	mux.Handle(cooperation_v1connect.NewCooperationServiceHandler(rpcServices.cooperation))
	mux.Handle(coordinator_v1connect.NewCoordinatorServiceHandler(rpcServices.coordinator))
	mux.Handle(event_v1connect.NewEventServiceHandler(rpcServices.event))
	mux.Handle(event_chat_v1connect.NewEventChatServiceHandler(rpcServices.eventChat))
	mux.Handle(geo_v1connect.NewGeoServiceHandler(rpcServices.geo))
	mux.Handle(participation_v1connect.NewParticipationServiceHandler(rpcServices.participation))
	mux.Handle(reward_v1connect.NewRewardServiceHandler(rpcServices.reward))
	mux.Handle(reward_order_v1connect.NewRewardOrderServiceHandler(rpcServices.rewardOrder))
	mux.Handle(volunteer_v1connect.NewVolunteerServiceHandler(rpcServices.volunteer))
	mux.Handle(session_v1connect.NewSessionServiceHandler(rpcServices.session))
	mux.Handle(image_v1connect.NewImageServiceHandler(rpcServices.image, connect.WithReadMaxBytes(imageReadLimit)))
	mux.Handle("/storage/", http.StripPrefix("/storage/", http.FileServer(http.Dir(cfg.Storage.Directory))))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })

	authentication := session_authentication_interceptor.NewAuthentication(
		supportServices.cookies,
		repos.sessions,
		tokens.sessionToken,
	)
	clientHandler := authn.NewMiddleware(authentication.Authenticate).Wrap(mux)
	clientHandler = withStreamingController(clientHandler)
	clientHandler = withCORS(clientHandler)

	return clientHandler, withCORS(withAdminAuthentication(adminMux, cfg.AdminAPIKey))
}
