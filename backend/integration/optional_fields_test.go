//go:build integration

package integration

import (
	"reflect"
	"time"
)

func checkOptionalFields(x *fixture) {
	const description = "Описание"

	owner := x.coordinator("optional-owner")
	x.call(owner, "coordinator.v1.CoordinatorService/UpdateCoordinator", map[string]any{"description": description})
	x.call(owner, "coordinator.v1.CoordinatorService/UpdateCoordinator", nil)
	coordinator := object(x.t, x.call(owner, "coordinator.v1.CoordinatorService/GetCoordinator", nil), "coordinator")
	require(x.t, coordinator["description"] == description, "omitted coordinator description was changed")
	x.call(owner, "coordinator.v1.CoordinatorService/UpdateCoordinator", map[string]any{"description": ""})
	coordinator = object(x.t, x.call(owner, "coordinator.v1.CoordinatorService/GetCoordinator", nil), "coordinator")
	require(x.t, coordinator["description"] == nil || coordinator["description"] == "", "description was not cleared")

	before := object(x.t, x.call(owner, "volunteer.v1.VolunteerService/GetVolunteer", nil), "volunteer")
	x.call(owner, "volunteer.v1.VolunteerService/UpdateVolunteer", map[string]any{"middle_name": "Отчество"})
	x.call(owner, "volunteer.v1.VolunteerService/UpdateVolunteer", map[string]any{"first_name": "Имя"})
	after := object(x.t, x.call(owner, "volunteer.v1.VolunteerService/GetVolunteer", nil), "volunteer")
	require(x.t, after["first_name"] == "Имя" && after["middle_name"] == "Отчество", "partial name update failed")
	require(x.t, after["last_name"] == before["last_name"], "omitted last name was changed")
	require(x.t, reflect.DeepEqual(after["date_of_birth"], before["date_of_birth"]), "omitted birthdate was changed")
	x.call(owner, "volunteer.v1.VolunteerService/UpdateVolunteer", map[string]any{"middle_name": ""})
	after = object(x.t, x.call(owner, "volunteer.v1.VolunteerService/GetVolunteer", nil), "volunteer")
	require(x.t, after["middle_name"] == nil || after["middle_name"] == "", "middle name was not cleared")

	eventID := x.event(owner, "Неизменное название", 2)
	before = object(x.t, x.call(owner, "event.v1.EventService/GetEvent", map[string]any{"id": eventID}), "event")
	x.call(
		owner,
		"event.v1.EventService/UpdateEvent",
		map[string]any{
			"id":          eventID,
			"description": description,
		},
	)

	after = object(x.t, x.call(owner, "event.v1.EventService/GetEvent", map[string]any{"id": eventID}), "event")
	for _, field := range []string{
		"title",
		"location",
		"point",
		"starts_at",
		"ends_at",
		"volunteer_limit",
	} {
		require(x.t, reflect.DeepEqual(after[field], before[field]), "omitted event field was changed: "+field)
	}

	x.call(owner, "event.v1.EventService/UpdateEvent", map[string]any{
		"id":          eventID,
		"description": "",
	})
	after = object(x.t, x.call(owner, "event.v1.EventService/GetEvent", map[string]any{"id": eventID}), "event")
	require(x.t, after["description"] == nil || after["description"] == "", "event description was not cleared")

	rewardID := text(x.t, x.call(owner, "reward.v1.RewardService/CreateReward", map[string]any{
		"title":          "Поощрение",
		"price_coins":    "10",
		"quantity_limit": 2,
	})["id"])
	x.call(
		owner,
		"reward.v1.RewardService/UpdateReward",
		map[string]any{
			"id":          rewardID,
			"description": description,
		},
	)
	x.call(owner, "reward.v1.RewardService/UpdateReward", map[string]any{
		"id":             rewardID,
		"quantity_limit": 0,
	})
	reward := object(x.t, x.call(owner, "reward.v1.RewardService/GetReward", map[string]any{"id": rewardID}), "reward")
	require(x.t, reward["title"] == "Поощрение" && integer(x.t, reward, "price_coins") == 10, "reward fields lost")
	require(x.t, reward["description"] == description, "omitted reward description was changed")
	require(x.t, integer(x.t, reward, "quantity_limit") == 0, "explicit zero quantity was ignored")
	x.fail(
		owner,
		"reward.v1.RewardService/UpdateReward",
		map[string]any{
			"id":          rewardID,
			"price_coins": "0",
		},
		"invalid_argument",
	)

	x.call(owner, "coordinator.v1.CoordinatorService/UpdateCoordinatorSettings", map[string]any{
		"cooperation_request_cooldown_seconds": "61",
	})
	x.call(owner, "coordinator.v1.CoordinatorService/UpdateCoordinatorSettings", nil)
	settings := object(x.t, x.call(owner, "coordinator.v1.CoordinatorService/GetCoordinatorSettings", nil), "settings")
	require(x.t, integer(x.t, settings, "cooperation_request_cooldown_seconds") == 61, "omitted cooldown was changed")
	x.call(owner, "coordinator.v1.CoordinatorService/UpdateCoordinatorSettings", map[string]any{
		"cooperation_request_cooldown_seconds": "0",
	})
	settings = object(x.t, x.call(owner, "coordinator.v1.CoordinatorService/GetCoordinatorSettings", nil), "settings")
	require(x.t, integer(x.t, settings, "cooperation_request_cooldown_seconds") == 0, "explicit zero cooldown ignored")

	expiresAt := time.Now().UTC().Add(time.Hour)
	created := x.call(owner, "cooperation.v1.CooperationService/CreateCooperationInvitation", map[string]any{
		"expires_at": expiresAt.Format(time.RFC3339Nano),
	})
	invitation := object(x.t, created, "invitation")
	actualExpiresAt, err := time.Parse(time.RFC3339Nano, text(x.t, invitation["expires_at"]))
	must(x.t, err)
	require(x.t, actualExpiresAt.Equal(expiresAt), "provided invitation expiration was ignored")
	x.fail(owner, "cooperation.v1.CooperationService/CreateCooperationInvitation", map[string]any{
		"expires_at": time.Now().UTC().Add(-time.Hour).Format(time.RFC3339Nano),
	}, "invalid_argument")

	x.call(owner, "event_chat.v1.EventChatService/SendEventChatMessage", map[string]any{
		"event_id":   eventID,
		"text":       "Сообщение",
		"request_id": "optional-snapshot",
	})
	messages := x.call(
		owner,
		"event_chat.v1.EventChatService/ListEventChatMessages",
		map[string]any{"event_id": eventID},
	)
	require(x.t, len(list(x.t, messages, "messages")) == 1, "omitted snapshot did not use latest sequence")
	messages = x.call(owner, "event_chat.v1.EventChatService/ListEventChatMessages", map[string]any{
		"event_id":          eventID,
		"snapshot_sequence": "0",
	})
	require(x.t, len(list(x.t, messages, "messages")) == 0, "explicit zero snapshot was ignored")
}
