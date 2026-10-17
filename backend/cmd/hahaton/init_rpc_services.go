package hahaton

import (
	rpc_service_admin "github.com/Roum1212/hahaton/internal/api/rpc/service/admin"
	rpc_service_admin_rpc_create_coordinator "github.com/Roum1212/hahaton/internal/api/rpc/service/admin/rpc/create-coordinator"
	rpc_service_admin_rpc_find_volunteers "github.com/Roum1212/hahaton/internal/api/rpc/service/admin/rpc/find-volunteers"
	rpc_service_cooperation "github.com/Roum1212/hahaton/internal/api/rpc/service/cooperation"
	rpc_service_cooperation_rpc_accept_cooperation_invitation "github.com/Roum1212/hahaton/internal/api/rpc/service/cooperation/rpc/accept-cooperation-invitation"
	rpc_service_cooperation_rpc_approve_cooperation "github.com/Roum1212/hahaton/internal/api/rpc/service/cooperation/rpc/approve-cooperation"
	rpc_service_cooperation_rpc_create_cooperation_invitation "github.com/Roum1212/hahaton/internal/api/rpc/service/cooperation/rpc/create-cooperation-invitation"
	rpc_service_cooperation_rpc_get_cooperation "github.com/Roum1212/hahaton/internal/api/rpc/service/cooperation/rpc/get-cooperation"
	rpc_service_cooperation_rpc_list_cooperation_invitations "github.com/Roum1212/hahaton/internal/api/rpc/service/cooperation/rpc/list-cooperation-invitations"
	rpc_service_cooperation_rpc_list_cooperations "github.com/Roum1212/hahaton/internal/api/rpc/service/cooperation/rpc/list-cooperations"
	rpc_service_cooperation_rpc_reject_cooperation "github.com/Roum1212/hahaton/internal/api/rpc/service/cooperation/rpc/reject-cooperation"
	rpc_service_cooperation_rpc_request_cooperation "github.com/Roum1212/hahaton/internal/api/rpc/service/cooperation/rpc/request-cooperation"
	rpc_service_cooperation_rpc_revoke_cooperation_invitation "github.com/Roum1212/hahaton/internal/api/rpc/service/cooperation/rpc/revoke-cooperation-invitation"
	rpc_service_cooperation_rpc_terminate_cooperation "github.com/Roum1212/hahaton/internal/api/rpc/service/cooperation/rpc/terminate-cooperation"
	rpc_service_coordinator "github.com/Roum1212/hahaton/internal/api/rpc/service/coordinator"
	rpc_service_coordinator_rpc_get_coordinator "github.com/Roum1212/hahaton/internal/api/rpc/service/coordinator/rpc/get-coordinator"
	rpc_service_coordinator_rpc_get_coordinator_settings "github.com/Roum1212/hahaton/internal/api/rpc/service/coordinator/rpc/get-coordinator-settings"
	rpc_service_coordinator_rpc_list_coordinators "github.com/Roum1212/hahaton/internal/api/rpc/service/coordinator/rpc/list-coordinators"
	rpc_service_coordinator_rpc_update_coordinator "github.com/Roum1212/hahaton/internal/api/rpc/service/coordinator/rpc/update-coordinator"
	rpc_service_coordinator_rpc_update_coordinator_settings "github.com/Roum1212/hahaton/internal/api/rpc/service/coordinator/rpc/update-coordinator-settings"
	rpc_service_event "github.com/Roum1212/hahaton/internal/api/rpc/service/event"
	rpc_service_event_chat "github.com/Roum1212/hahaton/internal/api/rpc/service/event-chat"
	rpc_service_event_chat_rpc_get_event_chat "github.com/Roum1212/hahaton/internal/api/rpc/service/event-chat/rpc/get-event-chat"
	rpc_service_event_chat_rpc_list_event_chat_messages "github.com/Roum1212/hahaton/internal/api/rpc/service/event-chat/rpc/list-event-chat-messages"
	rpc_service_event_chat_rpc_list_my_event_chats "github.com/Roum1212/hahaton/internal/api/rpc/service/event-chat/rpc/list-my-event-chats"
	rpc_service_event_chat_rpc_send_event_chat_message "github.com/Roum1212/hahaton/internal/api/rpc/service/event-chat/rpc/send-event-chat-message"
	rpc_service_event_chat_rpc_subscribe_event_chat_messages "github.com/Roum1212/hahaton/internal/api/rpc/service/event-chat/rpc/subscribe-event-chat-messages"
	rpc_service_event_rpc_cancel_event "github.com/Roum1212/hahaton/internal/api/rpc/service/event/rpc/cancel-event"
	rpc_service_event_rpc_complete_event "github.com/Roum1212/hahaton/internal/api/rpc/service/event/rpc/complete-event"
	rpc_service_event_rpc_create_event "github.com/Roum1212/hahaton/internal/api/rpc/service/event/rpc/create-event"
	rpc_service_event_rpc_get_event "github.com/Roum1212/hahaton/internal/api/rpc/service/event/rpc/get-event"
	rpc_service_event_rpc_list_events "github.com/Roum1212/hahaton/internal/api/rpc/service/event/rpc/list-events"
	rpc_service_event_rpc_update_event "github.com/Roum1212/hahaton/internal/api/rpc/service/event/rpc/update-event"
	rpc_service_geo "github.com/Roum1212/hahaton/internal/api/rpc/service/geo"
	rpc_service_geo_rpc_get_city "github.com/Roum1212/hahaton/internal/api/rpc/service/geo/rpc/get-city"
	rpc_service_geo_rpc_list_cities "github.com/Roum1212/hahaton/internal/api/rpc/service/geo/rpc/list-cities"
	rpc_service_geo_rpc_suggest_addresses "github.com/Roum1212/hahaton/internal/api/rpc/service/geo/rpc/suggest-addresses"
	rpc_service_image "github.com/Roum1212/hahaton/internal/api/rpc/service/image"
	rpc_service_image_rpc_upload_image "github.com/Roum1212/hahaton/internal/api/rpc/service/image/rpc/upload-image"
	rpc_service_participation "github.com/Roum1212/hahaton/internal/api/rpc/service/participation"
	rpc_service_participation_rpc_approve_participation "github.com/Roum1212/hahaton/internal/api/rpc/service/participation/rpc/approve-participation"
	rpc_service_participation_rpc_cancel_participation "github.com/Roum1212/hahaton/internal/api/rpc/service/participation/rpc/cancel-participation"
	rpc_service_participation_rpc_get_participation "github.com/Roum1212/hahaton/internal/api/rpc/service/participation/rpc/get-participation"
	rpc_service_participation_rpc_list_participation_changes "github.com/Roum1212/hahaton/internal/api/rpc/service/participation/rpc/list-participation-changes"
	rpc_service_participation_rpc_list_participations "github.com/Roum1212/hahaton/internal/api/rpc/service/participation/rpc/list-participations"
	rpc_service_participation_rpc_reject_participation "github.com/Roum1212/hahaton/internal/api/rpc/service/participation/rpc/reject-participation"
	rpc_service_participation_rpc_request_participation "github.com/Roum1212/hahaton/internal/api/rpc/service/participation/rpc/request-participation"
	rpc_service_participation_rpc_restore_participation "github.com/Roum1212/hahaton/internal/api/rpc/service/participation/rpc/restore-participation"
	rpc_service_reward "github.com/Roum1212/hahaton/internal/api/rpc/service/reward"
	rpc_service_reward_order "github.com/Roum1212/hahaton/internal/api/rpc/service/reward-order"
	rpc_service_reward_order_rpc_cancel_reward_order "github.com/Roum1212/hahaton/internal/api/rpc/service/reward-order/rpc/cancel-reward-order"
	rpc_service_reward_order_rpc_create_reward_order "github.com/Roum1212/hahaton/internal/api/rpc/service/reward-order/rpc/create-reward-order"
	rpc_service_reward_order_rpc_get_reward_order "github.com/Roum1212/hahaton/internal/api/rpc/service/reward-order/rpc/get-reward-order"
	rpc_service_reward_order_rpc_issue_reward_order "github.com/Roum1212/hahaton/internal/api/rpc/service/reward-order/rpc/issue-reward-order"
	rpc_service_reward_order_rpc_list_reward_orders "github.com/Roum1212/hahaton/internal/api/rpc/service/reward-order/rpc/list-reward-orders"
	rpc_service_reward_rpc_create_reward "github.com/Roum1212/hahaton/internal/api/rpc/service/reward/rpc/create-reward"
	rpc_service_reward_rpc_delete_reward "github.com/Roum1212/hahaton/internal/api/rpc/service/reward/rpc/delete-reward"
	rpc_service_reward_rpc_get_reward "github.com/Roum1212/hahaton/internal/api/rpc/service/reward/rpc/get-reward"
	rpc_service_reward_rpc_list_rewards "github.com/Roum1212/hahaton/internal/api/rpc/service/reward/rpc/list-rewards"
	rpc_service_reward_rpc_update_reward "github.com/Roum1212/hahaton/internal/api/rpc/service/reward/rpc/update-reward"
	rpc_service_session "github.com/Roum1212/hahaton/internal/api/rpc/service/session"
	rpc_service_session_rpc_authenticate_session "github.com/Roum1212/hahaton/internal/api/rpc/service/session/rpc/authenticate-session"
	rpc_service_session_rpc_close_session "github.com/Roum1212/hahaton/internal/api/rpc/service/session/rpc/close-session"
	rpc_service_session_rpc_get_session "github.com/Roum1212/hahaton/internal/api/rpc/service/session/rpc/get-session"
	rpc_service_session_rpc_open_session "github.com/Roum1212/hahaton/internal/api/rpc/service/session/rpc/open-session"
	rpc_service_volunteer "github.com/Roum1212/hahaton/internal/api/rpc/service/volunteer"
	rpc_service_volunteer_rpc_get_coordinator_volunteer "github.com/Roum1212/hahaton/internal/api/rpc/service/volunteer/rpc/get-coordinator-volunteer"
	rpc_service_volunteer_rpc_get_volunteer "github.com/Roum1212/hahaton/internal/api/rpc/service/volunteer/rpc/get-volunteer"
	rpc_service_volunteer_rpc_list_volunteer_operations "github.com/Roum1212/hahaton/internal/api/rpc/service/volunteer/rpc/list-volunteer-operations"
	rpc_service_volunteer_rpc_update_volunteer "github.com/Roum1212/hahaton/internal/api/rpc/service/volunteer/rpc/update-volunteer"
	"github.com/Roum1212/hahaton/internal/config"
)

type rpcServiceSet struct {
	admin         rpc_service_admin.Service
	cooperation   rpc_service_cooperation.Service
	coordinator   rpc_service_coordinator.Service
	event         rpc_service_event.Service
	eventChat     rpc_service_event_chat.Service
	geo           rpc_service_geo.Service
	image         rpc_service_image.Service
	participation rpc_service_participation.Service
	reward        rpc_service_reward.Service
	rewardOrder   rpc_service_reward_order.Service
	session       rpc_service_session.Service
	volunteer     rpc_service_volunteer.Service
}

func initRPCServices(cfg config.Config, services serviceSet, supportServices supportServiceSet) rpcServiceSet {
	return rpcServiceSet{
		admin: rpc_service_admin.NewService(
			rpc_service_admin_rpc_create_coordinator.NewRPC(services.createCoordinator),
			rpc_service_admin_rpc_find_volunteers.NewRPC(services.findVolunteers),
		),
		cooperation: rpc_service_cooperation.NewService(
			rpc_service_cooperation_rpc_accept_cooperation_invitation.NewRPC(services.acceptCooperationInvitation),
			rpc_service_cooperation_rpc_approve_cooperation.NewRPC(services.approveCooperation),
			rpc_service_cooperation_rpc_create_cooperation_invitation.NewRPC(
				services.createCooperationInvitation,
				cfg.FrontendURL,
			),
			rpc_service_cooperation_rpc_get_cooperation.NewRPC(services.getCooperation),
			rpc_service_cooperation_rpc_list_cooperation_invitations.NewRPC(
				services.listCooperationInvitations,
				cfg.FrontendURL,
			),
			rpc_service_cooperation_rpc_list_cooperations.NewRPC(services.listCooperations),
			rpc_service_cooperation_rpc_reject_cooperation.NewRPC(services.rejectCooperation),
			rpc_service_cooperation_rpc_request_cooperation.NewRPC(services.requestCooperation),
			rpc_service_cooperation_rpc_revoke_cooperation_invitation.NewRPC(services.revokeCooperationInvitation),
			rpc_service_cooperation_rpc_terminate_cooperation.NewRPC(services.terminateCooperation),
		),
		coordinator: rpc_service_coordinator.NewService(
			rpc_service_coordinator_rpc_get_coordinator.NewRPC(services.getCoordinator),
			rpc_service_coordinator_rpc_get_coordinator_settings.NewRPC(services.getCoordinatorSettings),
			rpc_service_coordinator_rpc_list_coordinators.NewRPC(services.listCoordinators),
			rpc_service_coordinator_rpc_update_coordinator.NewRPC(services.updateCoordinator),
			rpc_service_coordinator_rpc_update_coordinator_settings.NewRPC(services.updateCoordinatorSettings),
		),
		event: rpc_service_event.NewService(
			rpc_service_event_rpc_cancel_event.NewRPC(services.cancelEvent),
			rpc_service_event_rpc_complete_event.NewRPC(services.completeEvent),
			rpc_service_event_rpc_create_event.NewRPC(services.createEvent),
			rpc_service_event_rpc_get_event.NewRPC(services.getEvent),
			rpc_service_event_rpc_list_events.NewRPC(services.listEvents),
			rpc_service_event_rpc_update_event.NewRPC(services.updateEvent),
		),
		eventChat: rpc_service_event_chat.NewService(
			rpc_service_event_chat_rpc_get_event_chat.NewRPC(services.getEventChat),
			rpc_service_event_chat_rpc_list_event_chat_messages.NewRPC(services.listEventChatMessages),
			rpc_service_event_chat_rpc_list_my_event_chats.NewRPC(services.listMyEventChats),
			rpc_service_event_chat_rpc_send_event_chat_message.NewRPC(services.sendEventChatMessage),
			rpc_service_event_chat_rpc_subscribe_event_chat_messages.NewRPC(services.subscribeEventChatMessages),
		),
		geo: rpc_service_geo.NewService(
			rpc_service_geo_rpc_get_city.NewRPC(services.getCity),
			rpc_service_geo_rpc_list_cities.NewRPC(services.listCities),
			rpc_service_geo_rpc_suggest_addresses.NewRPC(services.suggestAddresses),
		),
		image: rpc_service_image.NewService(rpc_service_image_rpc_upload_image.NewRPC(services.uploadImage)),
		participation: rpc_service_participation.NewService(
			rpc_service_participation_rpc_approve_participation.NewRPC(services.approveParticipation),
			rpc_service_participation_rpc_cancel_participation.NewRPC(services.cancelParticipation),
			rpc_service_participation_rpc_get_participation.NewRPC(services.getParticipation),
			rpc_service_participation_rpc_list_participation_changes.NewRPC(services.listParticipationChanges),
			rpc_service_participation_rpc_list_participations.NewRPC(services.listParticipations),
			rpc_service_participation_rpc_reject_participation.NewRPC(services.rejectParticipation),
			rpc_service_participation_rpc_request_participation.NewRPC(services.requestParticipation),
			rpc_service_participation_rpc_restore_participation.NewRPC(services.restoreParticipation),
		),
		reward: rpc_service_reward.NewService(
			rpc_service_reward_rpc_create_reward.NewRPC(services.createReward),
			rpc_service_reward_rpc_delete_reward.NewRPC(services.deleteReward),
			rpc_service_reward_rpc_get_reward.NewRPC(services.getReward),
			rpc_service_reward_rpc_list_rewards.NewRPC(services.listRewards),
			rpc_service_reward_rpc_update_reward.NewRPC(services.updateReward),
		),
		rewardOrder: rpc_service_reward_order.NewService(
			rpc_service_reward_order_rpc_cancel_reward_order.NewRPC(services.cancelRewardOrder),
			rpc_service_reward_order_rpc_create_reward_order.NewRPC(services.createRewardOrder),
			rpc_service_reward_order_rpc_get_reward_order.NewRPC(services.getRewardOrder),
			rpc_service_reward_order_rpc_issue_reward_order.NewRPC(services.issueRewardOrder),
			rpc_service_reward_order_rpc_list_reward_orders.NewRPC(services.listRewardOrders),
		),
		session: rpc_service_session.NewService(
			rpc_service_session_rpc_authenticate_session.NewRPC(
				supportServices.cookies,
				services.getOrRegisterVolunteer,
				services.authenticateSession,
			),
			rpc_service_session_rpc_close_session.NewRPC(supportServices.expiredCookies, services.closeSession),
			rpc_service_session_rpc_get_session.NewRPC(services.getSession),
			rpc_service_session_rpc_open_session.NewRPC(supportServices.cookies, services.openSession)),
		volunteer: rpc_service_volunteer.NewService(
			rpc_service_volunteer_rpc_get_coordinator_volunteer.NewRPC(services.getCoordinatorVolunteer),
			rpc_service_volunteer_rpc_get_volunteer.NewRPC(services.getVolunteer),
			rpc_service_volunteer_rpc_list_volunteer_operations.NewRPC(services.listVolunteerOperations),
			rpc_service_volunteer_rpc_update_volunteer.NewRPC(services.updateVolunteer),
		),
	}
}
