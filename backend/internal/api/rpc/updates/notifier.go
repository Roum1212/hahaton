package rpc_updates

import (
	"context"

	"connectrpc.com/connect"

	admin_v1connect "github.com/Roum1212/hahaton/genproto/admin/v1/admin_v1connect"
	cooperation_v1connect "github.com/Roum1212/hahaton/genproto/cooperation/v1/cooperation_v1connect"
	coordinator_v1connect "github.com/Roum1212/hahaton/genproto/coordinator/v1/coordinator_v1connect"
	event_chat_v1connect "github.com/Roum1212/hahaton/genproto/event-chat/v1/event_chat_v1connect"
	event_v1connect "github.com/Roum1212/hahaton/genproto/event/v1/event_v1connect"
	participation_v1connect "github.com/Roum1212/hahaton/genproto/participation/v1/participation_v1connect"
	reward_order_v1connect "github.com/Roum1212/hahaton/genproto/reward-order/v1/reward_order_v1connect"
	reward_v1connect "github.com/Roum1212/hahaton/genproto/reward/v1/reward_v1connect"
	session_v1connect "github.com/Roum1212/hahaton/genproto/session/v1/session_v1connect"
	volunteer_v1connect "github.com/Roum1212/hahaton/genproto/volunteer/v1/volunteer_v1connect"
	rpc_updates_serde "github.com/Roum1212/hahaton/internal/api/rpc/updates/serde"
	service_get_update_recipients "github.com/Roum1212/hahaton/internal/app/service/get-update-recipients"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
	"github.com/Roum1212/hahaton/internal/pkg/updates"
)

type Notifier struct {
	hub        *updates.Hub
	recipients service_get_update_recipients.Service
}

func (x Notifier) Notify(ctx context.Context, r connect.AnyRequest, response connect.AnyResponse) {
	if !isMutation(r.Spec().Procedure) {
		return
	}

	var userID model_user_id.UserID
	if session, ok := session_context.GetSession(ctx); ok {
		userID = session.GetUserID()
	}

	request, err := rpc_updates_serde.RequestFromRPC(userID, r.Any(), response.Any())
	if err != nil {
		x.hub.NotifyAll()

		return
	}

	responseRecipients, err := x.recipients.GetUpdateRecipients(ctx, request)
	if err != nil {
		// The write has committed. A payload-free resync preserves its successful response.
		x.hub.NotifyAll()

		return
	}

	x.hub.Notify(responseRecipients.GetRecipients())
}

func isMutation(procedure string) bool {
	switch procedure {
	case admin_v1connect.AdminServiceCreateCoordinatorProcedure:
		return true
	case cooperation_v1connect.CooperationServiceRequestCooperationProcedure:
		return true
	case cooperation_v1connect.CooperationServiceApproveCooperationProcedure:
		return true
	case cooperation_v1connect.CooperationServiceRejectCooperationProcedure:
		return true
	case cooperation_v1connect.CooperationServiceTerminateCooperationProcedure:
		return true
	case cooperation_v1connect.CooperationServiceCreateCooperationInvitationProcedure:
		return true
	case cooperation_v1connect.CooperationServiceRevokeCooperationInvitationProcedure:
		return true
	case cooperation_v1connect.CooperationServiceAcceptCooperationInvitationProcedure:
		return true
	case coordinator_v1connect.CoordinatorServiceUpdateCoordinatorProcedure:
		return true
	case coordinator_v1connect.CoordinatorServiceUpdateCoordinatorSettingsProcedure:
		return true
	case event_v1connect.EventServiceCreateEventProcedure:
		return true
	case event_v1connect.EventServiceUpdateEventProcedure:
		return true
	case event_v1connect.EventServiceCancelEventProcedure:
		return true
	case event_v1connect.EventServiceCompleteEventProcedure:
		return true
	case event_chat_v1connect.EventChatServiceSendEventChatMessageProcedure:
		return true
	case participation_v1connect.ParticipationServiceRequestParticipationProcedure:
		return true
	case participation_v1connect.ParticipationServiceRestoreParticipationProcedure:
		return true
	case participation_v1connect.ParticipationServiceApproveParticipationProcedure:
		return true
	case participation_v1connect.ParticipationServiceRejectParticipationProcedure:
		return true
	case participation_v1connect.ParticipationServiceCancelParticipationProcedure:
		return true
	case reward_v1connect.RewardServiceCreateRewardProcedure:
		return true
	case reward_v1connect.RewardServiceUpdateRewardProcedure:
		return true
	case reward_v1connect.RewardServiceDeleteRewardProcedure:
		return true
	case reward_order_v1connect.RewardOrderServiceCreateRewardOrderProcedure:
		return true
	case reward_order_v1connect.RewardOrderServiceIssueRewardOrderProcedure:
		return true
	case reward_order_v1connect.RewardOrderServiceCancelRewardOrderProcedure:
		return true
	case session_v1connect.SessionServiceCloseSessionProcedure:
		return true
	case volunteer_v1connect.VolunteerServiceUpdateVolunteerProcedure:
		return true
	default:
		return false
	}
}

func NewNotifier(hub *updates.Hub, recipients service_get_update_recipients.Service) Notifier {
	return Notifier{hub: hub, recipients: recipients}
}
