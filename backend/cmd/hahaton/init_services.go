package hahaton

import (
	"net/http"

	service_accept_cooperation_invitation "github.com/Roum1212/hahaton/internal/app/service/accept-cooperation-invitation"
	service_approve_cooperation "github.com/Roum1212/hahaton/internal/app/service/approve-cooperation"
	service_approve_participation "github.com/Roum1212/hahaton/internal/app/service/approve-participation"
	service_authenticate_session "github.com/Roum1212/hahaton/internal/app/service/authenticate-session"
	service_cancel_event "github.com/Roum1212/hahaton/internal/app/service/cancel-event"
	service_cancel_participation "github.com/Roum1212/hahaton/internal/app/service/cancel-participation"
	service_cancel_reward_order "github.com/Roum1212/hahaton/internal/app/service/cancel-reward-order"
	service_close_session "github.com/Roum1212/hahaton/internal/app/service/close-session"
	service_complete_event "github.com/Roum1212/hahaton/internal/app/service/complete-event"
	service_create_cooperation_invitation "github.com/Roum1212/hahaton/internal/app/service/create-cooperation-invitation"
	service_create_coordinator "github.com/Roum1212/hahaton/internal/app/service/create-coordinator"
	service_create_event "github.com/Roum1212/hahaton/internal/app/service/create-event"
	service_create_reward "github.com/Roum1212/hahaton/internal/app/service/create-reward"
	service_create_reward_order "github.com/Roum1212/hahaton/internal/app/service/create-reward-order"
	service_delete_reward "github.com/Roum1212/hahaton/internal/app/service/delete-reward"
	service_find_volunteers "github.com/Roum1212/hahaton/internal/app/service/find-volunteers"
	service_get_city "github.com/Roum1212/hahaton/internal/app/service/get-city"
	service_get_cooperation "github.com/Roum1212/hahaton/internal/app/service/get-cooperation"
	service_get_coordinator "github.com/Roum1212/hahaton/internal/app/service/get-coordinator"
	service_get_coordinator_settings "github.com/Roum1212/hahaton/internal/app/service/get-coordinator-settings"
	service_get_coordinator_volunteer "github.com/Roum1212/hahaton/internal/app/service/get-coordinator-volunteer"
	service_get_event "github.com/Roum1212/hahaton/internal/app/service/get-event"
	service_get_event_chat "github.com/Roum1212/hahaton/internal/app/service/get-event-chat"
	service_get_or_register_volunteer "github.com/Roum1212/hahaton/internal/app/service/get-or-register-volunteer"
	service_get_participation "github.com/Roum1212/hahaton/internal/app/service/get-participation"
	service_get_reward "github.com/Roum1212/hahaton/internal/app/service/get-reward"
	service_get_reward_order "github.com/Roum1212/hahaton/internal/app/service/get-reward-order"
	service_get_session "github.com/Roum1212/hahaton/internal/app/service/get-session"
	service_get_update_recipients "github.com/Roum1212/hahaton/internal/app/service/get-update-recipients"
	service_get_volunteer "github.com/Roum1212/hahaton/internal/app/service/get-volunteer"
	service_issue_reward_order "github.com/Roum1212/hahaton/internal/app/service/issue-reward-order"
	service_list_cities "github.com/Roum1212/hahaton/internal/app/service/list-cities"
	service_list_cooperation_invitations "github.com/Roum1212/hahaton/internal/app/service/list-cooperation-invitations"
	service_list_cooperations "github.com/Roum1212/hahaton/internal/app/service/list-cooperations"
	service_list_coordinators "github.com/Roum1212/hahaton/internal/app/service/list-coordinators"
	service_list_event_chat_messages "github.com/Roum1212/hahaton/internal/app/service/list-event-chat-messages"
	service_list_events "github.com/Roum1212/hahaton/internal/app/service/list-events"
	service_list_my_event_chats "github.com/Roum1212/hahaton/internal/app/service/list-my-event-chats"
	service_list_participation_changes "github.com/Roum1212/hahaton/internal/app/service/list-participation-changes"
	service_list_participations "github.com/Roum1212/hahaton/internal/app/service/list-participations"
	service_list_reward_orders "github.com/Roum1212/hahaton/internal/app/service/list-reward-orders"
	service_list_rewards "github.com/Roum1212/hahaton/internal/app/service/list-rewards"
	service_list_volunteer_operations "github.com/Roum1212/hahaton/internal/app/service/list-volunteer-operations"
	service_open_session "github.com/Roum1212/hahaton/internal/app/service/open-session"
	service_reject_cooperation "github.com/Roum1212/hahaton/internal/app/service/reject-cooperation"
	service_reject_participation "github.com/Roum1212/hahaton/internal/app/service/reject-participation"
	service_request_cooperation "github.com/Roum1212/hahaton/internal/app/service/request-cooperation"
	service_request_participation "github.com/Roum1212/hahaton/internal/app/service/request-participation"
	service_restore_participation "github.com/Roum1212/hahaton/internal/app/service/restore-participation"
	service_revoke_cooperation_invitation "github.com/Roum1212/hahaton/internal/app/service/revoke-cooperation-invitation"
	service_send_event_chat_message "github.com/Roum1212/hahaton/internal/app/service/send-event-chat-message"
	service_subscribe_event_chat_messages "github.com/Roum1212/hahaton/internal/app/service/subscribe-event-chat-messages"
	service_subscribe_updates "github.com/Roum1212/hahaton/internal/app/service/subscribe-updates"
	service_suggest_addresses "github.com/Roum1212/hahaton/internal/app/service/suggest-addresses"
	service_terminate_cooperation "github.com/Roum1212/hahaton/internal/app/service/terminate-cooperation"
	service_update_coordinator "github.com/Roum1212/hahaton/internal/app/service/update-coordinator"
	service_update_coordinator_settings "github.com/Roum1212/hahaton/internal/app/service/update-coordinator-settings"
	service_update_event "github.com/Roum1212/hahaton/internal/app/service/update-event"
	service_update_reward "github.com/Roum1212/hahaton/internal/app/service/update-reward"
	service_update_volunteer "github.com/Roum1212/hahaton/internal/app/service/update-volunteer"
	service_upload_image "github.com/Roum1212/hahaton/internal/app/service/upload-image"
	"github.com/Roum1212/hahaton/internal/config"
	"github.com/Roum1212/hahaton/internal/pkg/cookies"
	session_cookie "github.com/Roum1212/hahaton/internal/pkg/session/cookie"
)

type serviceSet struct {
	acceptCooperationInvitation service_accept_cooperation_invitation.Service
	approveCooperation          service_approve_cooperation.Service
	approveParticipation        service_approve_participation.Service
	authenticateSession         service_authenticate_session.Service
	cancelEvent                 service_cancel_event.Service
	cancelParticipation         service_cancel_participation.Service
	cancelRewardOrder           service_cancel_reward_order.Service
	closeSession                service_close_session.Service
	completeEvent               service_complete_event.Service
	cookie                      cookies.Service
	createCooperationInvitation service_create_cooperation_invitation.Service
	createCoordinator           service_create_coordinator.Service
	createEvent                 service_create_event.Service
	createReward                service_create_reward.Service
	createRewardOrder           service_create_reward_order.Service
	deleteReward                service_delete_reward.Service
	findVolunteers              service_find_volunteers.Service
	getCity                     service_get_city.Service
	getCooperation              service_get_cooperation.Service
	getCoordinator              service_get_coordinator.Service
	getCoordinatorSettings      service_get_coordinator_settings.Service
	getCoordinatorVolunteer     service_get_coordinator_volunteer.Service
	getEvent                    service_get_event.Service
	getEventChat                service_get_event_chat.Service
	getOrRegisterVolunteer      service_get_or_register_volunteer.Service
	getParticipation            service_get_participation.Service
	getReward                   service_get_reward.Service
	getRewardOrder              service_get_reward_order.Service
	getSession                  service_get_session.Service
	getUpdateRecipients         service_get_update_recipients.Service
	getVolunteer                service_get_volunteer.Service
	issueRewardOrder            service_issue_reward_order.Service
	listCities                  service_list_cities.Service
	listCooperationInvitations  service_list_cooperation_invitations.Service
	listCooperations            service_list_cooperations.Service
	listCoordinators            service_list_coordinators.Service
	listEventChatMessages       service_list_event_chat_messages.Service
	listEvents                  service_list_events.Service
	listMyEventChats            service_list_my_event_chats.Service
	listParticipationChanges    service_list_participation_changes.Service
	listParticipations          service_list_participations.Service
	listRewardOrders            service_list_reward_orders.Service
	listRewards                 service_list_rewards.Service
	listVolunteerOperations     service_list_volunteer_operations.Service
	openSession                 service_open_session.Service
	rejectCooperation           service_reject_cooperation.Service
	rejectParticipation         service_reject_participation.Service
	requestCooperation          service_request_cooperation.Service
	requestParticipation        service_request_participation.Service
	restoreParticipation        service_restore_participation.Service
	revokeCooperationInvitation service_revoke_cooperation_invitation.Service
	sendEventChatMessage        service_send_event_chat_message.Service
	subscribeEventChatMessages  service_subscribe_event_chat_messages.Service
	subscribeUpdates            service_subscribe_updates.Service
	suggestAddresses            service_suggest_addresses.Service
	terminateCooperation        service_terminate_cooperation.Service
	updateCoordinator           service_update_coordinator.Service
	updateCoordinatorSettings   service_update_coordinator_settings.Service
	updateEvent                 service_update_event.Service
	updateReward                service_update_reward.Service
	updateVolunteer             service_update_volunteer.Service
	uploadImage                 service_upload_image.Service
}

func initServices(
	cfg config.Config,
	externalClients externalClientSet,
	repos repositorySet,
	tokens tokenSet,
) serviceSet {
	return serviceSet{
		acceptCooperationInvitation: service_accept_cooperation_invitation.NewService(
			repos.invitations,
			repos.cooperations,
		),
		approveCooperation: service_approve_cooperation.NewService(repos.cooperations, repos.coordinators),
		approveParticipation: service_approve_participation.NewService(
			repos.events,
			repos.coordinators,
			repos.participations,
			repos.cooperations,
			repos.chats,
		),
		authenticateSession: service_authenticate_session.NewService(
			repos.sessions,
			tokens.sessionToken,
			cfg.SessionToken.TTL,
		),
		cancelEvent:         service_cancel_event.NewService(repos.events, repos.coordinators, repos.participations),
		cancelParticipation: service_cancel_participation.NewService(repos.events, repos.participations),
		cancelRewardOrder: service_cancel_reward_order.NewService(
			repos.orders,
			repos.coordinators,
			repos.rewards,
			repos.wallets,
		),
		closeSession: service_close_session.NewService(repos.sessions),
		completeEvent: service_complete_event.NewService(
			repos.events,
			repos.coordinators,
			repos.participations,
			repos.wallets,
		),
		cookie: cookies.NewService(
			session_cookie.Name,
			cookies.WithHTTPOnly(),
			cookies.WithSecure(),
			cookies.WithPath("/"),
			cookies.WithSameSite(http.SameSiteLaxMode),
			cookies.WithMaxAge(cfg.SessionToken.TTL),
		),
		createCooperationInvitation: service_create_cooperation_invitation.NewService(
			repos.invitations,
			repos.coordinators,
		),
		createCoordinator: service_create_coordinator.NewService(
			repos.volunteers,
			repos.coordinators,
			externalClients.images,
		),
		createEvent: service_create_event.NewService(
			repos.events,
			repos.coordinators,
			repos.chats,
			externalClients.images,
		),
		createReward: service_create_reward.NewService(repos.rewards, repos.coordinators, externalClients.images),
		createRewardOrder: service_create_reward_order.NewService(
			repos.rewards,
			repos.cooperations,
			repos.orders,
			repos.wallets,
		),
		deleteReward: service_delete_reward.NewService(
			repos.rewards,
			repos.coordinators,
			repos.orders,
			repos.wallets,
		),
		findVolunteers:         service_find_volunteers.NewService(repos.volunteers, repos.coordinators),
		getCity:                service_get_city.NewService(externalClients.dadata),
		getCooperation:         service_get_cooperation.NewService(repos.cooperations, repos.coordinators),
		getCoordinator:         service_get_coordinator.NewService(repos.coordinators, repos.volunteers),
		getCoordinatorSettings: service_get_coordinator_settings.NewService(repos.coordinators),
		getCoordinatorVolunteer: service_get_coordinator_volunteer.NewService(
			repos.coordinators,
			repos.volunteers,
			repos.cooperations,
			repos.participations,
			repos.orders,
		),
		getEvent: service_get_event.NewService(repos.events, repos.participations, repos.cooperations),
		getEventChat: service_get_event_chat.NewService(
			repos.events,
			repos.coordinators,
			repos.cooperations,
			repos.participations,
			repos.chats,
		),
		getOrRegisterVolunteer: service_get_or_register_volunteer.NewService(
			externalClients.flida,
			repos.volunteers,
			repos.wallets,
		),
		getParticipation: service_get_participation.NewService(repos.participations, repos.coordinators),
		getReward:        service_get_reward.NewService(repos.rewards, repos.cooperations),
		getRewardOrder:   service_get_reward_order.NewService(repos.orders, repos.coordinators),
		getSession:       service_get_session.NewService(repos.coordinators),
		getUpdateRecipients: service_get_update_recipients.NewService(
			repos.cooperations, repos.events, repos.participations, repos.rewards, repos.orders,
		),
		getVolunteer: service_get_volunteer.NewService(repos.volunteers, repos.wallets),
		issueRewardOrder: service_issue_reward_order.NewService(
			repos.orders,
			repos.coordinators,
			repos.rewards,
			repos.cooperations,
		),
		listCities: service_list_cities.NewService(externalClients.dadata),
		listCooperationInvitations: service_list_cooperation_invitations.NewService(
			repos.invitations,
			repos.coordinators,
		),
		listCooperations: service_list_cooperations.NewService(repos.cooperations, repos.coordinators),
		listCoordinators: service_list_coordinators.NewService(repos.coordinators, repos.volunteers),
		listEventChatMessages: service_list_event_chat_messages.NewService(
			repos.events,
			repos.coordinators,
			repos.cooperations,
			repos.participations,
			repos.chats,
		),
		listEvents: service_list_events.NewService(
			repos.events,
			repos.coordinators,
			repos.cooperations,
			repos.participations,
			externalClients.dadata,
		),
		listMyEventChats: service_list_my_event_chats.NewService(
			repos.events,
			repos.coordinators,
			repos.cooperations,
			repos.participations,
			repos.chats,
		),
		listParticipationChanges: service_list_participation_changes.NewService(
			repos.participations,
			repos.coordinators,
		),
		listParticipations:      service_list_participations.NewService(repos.participations, repos.coordinators),
		listRewardOrders:        service_list_reward_orders.NewService(repos.orders, repos.coordinators),
		listRewards:             service_list_rewards.NewService(repos.rewards, repos.coordinators, repos.cooperations),
		listVolunteerOperations: service_list_volunteer_operations.NewService(repos.wallets),
		openSession: service_open_session.NewService(
			repos.sessions,
			tokens.sessionToken,
			cfg.SessionToken.TTL,
		),
		rejectCooperation: service_reject_cooperation.NewService(repos.cooperations, repos.coordinators),
		rejectParticipation: service_reject_participation.NewService(
			repos.events,
			repos.coordinators,
			repos.participations,
		),
		requestCooperation: service_request_cooperation.NewService(repos.cooperations, repos.coordinators),
		requestParticipation: service_request_participation.NewService(
			repos.events,
			repos.participations,
			repos.cooperations,
		),
		restoreParticipation: service_restore_participation.NewService(
			repos.events,
			repos.participations,
			repos.cooperations,
		),
		revokeCooperationInvitation: service_revoke_cooperation_invitation.NewService(
			repos.invitations,
			repos.coordinators,
		),
		sendEventChatMessage: service_send_event_chat_message.NewService(
			repos.events,
			repos.coordinators,
			repos.cooperations,
			repos.participations,
			repos.chats,
			repos.volunteers,
		),
		subscribeEventChatMessages: service_subscribe_event_chat_messages.NewService(
			repos.events,
			repos.coordinators,
			repos.cooperations,
			repos.participations,
			repos.chats,
			repos.sessions,
		),
		subscribeUpdates: service_subscribe_updates.NewService(repos.sessions),
		suggestAddresses: service_suggest_addresses.NewService(externalClients.dadata),
		terminateCooperation: service_terminate_cooperation.NewService(
			repos.cooperations,
			repos.coordinators,
			repos.participations,
			repos.orders,
			repos.rewards,
			repos.wallets,
		),
		updateCoordinator:         service_update_coordinator.NewService(repos.coordinators, externalClients.images),
		updateCoordinatorSettings: service_update_coordinator_settings.NewService(repos.coordinators),
		updateEvent: service_update_event.NewService(
			repos.events,
			repos.coordinators,
			repos.participations,
			externalClients.images,
		),
		updateReward:    service_update_reward.NewService(repos.rewards, repos.coordinators, externalClients.images),
		updateVolunteer: service_update_volunteer.NewService(repos.volunteers, externalClients.images),
		uploadImage:     service_upload_image.NewService(externalClients.images),
	}
}
