//go:build integration

package integration

import (
	"encoding/base64"
	"testing"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/proto"
)

func badRequestField(t *testing.T, result map[string]any, field string) {
	t.Helper()

	details := list(t, result, "details")
	require(t, len(details) == 1, "invalid argument must include bad request details")
	detail := dictionary(t, details[0])
	require(t, detail["type"] == "google.rpc.BadRequest", "unexpected invalid argument detail type")
	data, err := base64.RawStdEncoding.DecodeString(text(t, detail["value"]))
	must(t, err)

	info := &errdetails.BadRequest{}
	must(t, proto.Unmarshal(data, info))
	violations := info.GetFieldViolations()
	require(t, len(violations) == 1, "expected one field violation")
	require(t, violations[0].GetField() == field, "wrong field violation path: "+violations[0].GetField())
	require(t, violations[0].GetDescription() != "", "field violation description missing")
}

func checkAPIConversions(x *fixture) {
	c := x.coordinator("api-conversions")

	for _, method := range []string{
		"event.v1.EventService/GetEvent",
		"participation.v1.ParticipationService/GetParticipation",
		"reward.v1.RewardService/GetReward",
		"reward_order.v1.RewardOrderService/GetRewardOrder",
		"cooperation.v1.CooperationService/RevokeCooperationInvitation",
	} {
		result := x.fail(c, method, map[string]any{"id": "invalid"}, "invalid_argument")
		badRequestField(x.t, result, "id")
	}

	for _, field := range []string{"counterpart", "coordinator_id", "volunteer_id"} {
		body := make(map[string]any)
		if field != "counterpart" {
			body[field] = "invalid"
		}

		result := x.fail(c, "cooperation.v1.CooperationService/GetCooperation", body, "invalid_argument")
		badRequestField(x.t, result, field)
	}

	id := x.event(c, "Преобразования API", 1)
	for _, requirement := range []string{
		"EVENT_CANDIDATE_REQUIREMENT_AGE_14_PLUS",
		"EVENT_CANDIDATE_REQUIREMENT_AGE_16_PLUS",
		"EVENT_CANDIDATE_REQUIREMENT_AGE_18_PLUS",
		"EVENT_CANDIDATE_REQUIREMENT_PRELIMINARY_TRAINING",
		"EVENT_CANDIDATE_REQUIREMENT_RUSSIAN_CITIZENSHIP",
		"EVENT_CANDIDATE_REQUIREMENT_VOLUNTEERING_EXPERIENCE",
	} {
		x.call(c, "event.v1.EventService/UpdateEvent", map[string]any{
			"id":                     id,
			"candidate_requirements": map[string]any{"values": []string{requirement}},
		})
		event := object(x.t, x.call(c, "event.v1.EventService/GetEvent", map[string]any{"id": id}), "event")
		values := list(x.t, event, "candidate_requirements")
		require(x.t, len(values) == 1 && values[0] == requirement, "requirement changed during round trip")
	}

	for _, values := range [][]int{{0}, {99}, {1, 1}} {
		result := x.fail(c, "event.v1.EventService/UpdateEvent", map[string]any{
			"id":                     id,
			"candidate_requirements": map[string]any{"values": values},
		}, "invalid_argument")
		badRequestField(x.t, result, "candidate_requirements")
	}

	for _, operationType := range []int{0, 99} {
		result := x.fail(c, "volunteer.v1.VolunteerService/ListVolunteerOperations", map[string]any{
			"type": operationType,
		}, "invalid_argument")
		badRequestField(x.t, result, "type")
	}

	for _, operationType := range []int{1, 2, 3} {
		x.call(c, "volunteer.v1.VolunteerService/ListVolunteerOperations", map[string]any{"type": operationType})
	}

	for _, outcome := range []map[string]any{{}, {"absent": false}, {"confirmed_minutes": "0"}, {
		"confirmed_minutes": "45",
	}} {
		outcome["participation_id"] = c.id
		result := x.fail(c, "event.v1.EventService/CompleteEvent", map[string]any{
			"id":      id,
			"results": []map[string]any{outcome},
		}, "invalid_argument")
		badRequestField(x.t, result, "results[0]")
	}
}
