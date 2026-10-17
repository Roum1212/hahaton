//go:build integration

package integration

import (
	"fmt"
	"time"
)

func checkEvents(x *fixture) {
	c, d := x.coordinator("events-owner"), x.coordinator("events-other")

	a, b, waiting := x.login("event-alice"), x.login("event-bob"), x.login("event-waiting")
	for _, volunteer := range []actor{
		a,
		b,
		waiting,
	} {
		x.cooperate(c, volunteer)
	}

	id := x.event(c, "Акция", 1)
	x.event(c, "Вторая акция", 2)
	x.event(c, "Третья акция", 2)
	own := x.call(c, "event.v1.EventService/ListEvents", map[string]any{
		"scope": 2,
		"limit": 1,
	})
	require(
		x.t,
		len(list(x.t, own, "events")) == 1 && fmt.Sprint(own["total_count"]) == "3",
		"event pagination/count differ",
	)
	x.fail(a, "event.v1.EventService/ListEvents", map[string]any{"scope": 2}, "permission_denied")
	x.fail(a, "event.v1.EventService/ListEvents", map[string]any{"scope": 1}, "invalid_argument")
	x.fail(d, "event.v1.EventService/GetEvent", map[string]any{"id": id}, "permission_denied")
	x.fail(d, "event.v1.EventService/UpdateEvent", map[string]any{
		"id":    id,
		"title": "Чужое",
	}, "permission_denied")
	x.call(c, "event.v1.EventService/UpdateEvent", map[string]any{
		"id":          id,
		"description": "Описание",
		"candidate_requirements": map[string]any{
			"values": []int{
				3,
				4,
			},
		},
	})
	x.call(c, "event.v1.EventService/UpdateEvent", map[string]any{
		"id":                     id,
		"candidate_requirements": map[string]any{},
	})
	first := text(x.t, x.call(a, "participation.v1.ParticipationService/RequestParticipation", map[string]any{
		"event_id": id,
	})["id"])
	x.fail(a, "participation.v1.ParticipationService/RequestParticipation", map[string]any{
		"event_id": id,
	}, "already_exists")
	x.fail(c, "event.v1.EventService/UpdateEvent", map[string]any{
		"id":          id,
		"description": "Запрещено",
	}, "failed_precondition")
	x.call(c, "participation.v1.ParticipationService/ApproveParticipation", map[string]any{"id": first})
	second := text(x.t, x.call(b, "participation.v1.ParticipationService/RequestParticipation", map[string]any{
		"event_id": id,
	})["id"])
	x.fail(c, "participation.v1.ParticipationService/ApproveParticipation", map[string]any{
		"id": second,
	}, "failed_precondition")
	x.call(a, "participation.v1.ParticipationService/CancelParticipation", map[string]any{"id": first})
	x.call(c, "participation.v1.ParticipationService/ApproveParticipation", map[string]any{"id": second})
	x.call(a, "participation.v1.ParticipationService/RestoreParticipation", map[string]any{"id": first})
	x.call(c, "participation.v1.ParticipationService/RejectParticipation", map[string]any{
		"id":     second,
		"reason": "  Условия изменились  ",
	})
	rejection := object(x.t, x.call(b, "participation.v1.ParticipationService/GetParticipation", map[string]any{
		"id": second,
	}), "participation")
	require(x.t, rejection["rejection_reason"] == "Условия изменились", "rejection reason not normalized")
	x.fail(a, "participation.v1.ParticipationService/GetParticipation", map[string]any{
		"id": second,
	}, "permission_denied")
	x.fail(b, "participation.v1.ParticipationService/RestoreParticipation", map[string]any{
		"id": second,
	}, "failed_precondition")
	x.fail(b, "participation.v1.ParticipationService/RequestParticipation", map[string]any{
		"event_id": id,
	}, "already_exists")
	x.call(c, "participation.v1.ParticipationService/ApproveParticipation", map[string]any{"id": first})
	pending := text(x.t, x.call(waiting, "participation.v1.ParticipationService/RequestParticipation", map[string]any{
		"event_id": id,
	})["id"])
	x.fail(c, "event.v1.EventService/CompleteEvent", map[string]any{
		"id":         id,
		"request_id": "complete",
		"results": []any{
			map[string]any{
				"participation_id":  first,
				"confirmed_minutes": "180",
			},
		},
	}, "failed_precondition")
	before := object(x.t, x.call(a, "volunteer.v1.VolunteerService/GetVolunteer", nil), "volunteer")
	require(x.t, integer(x.t, before, "coin_balance") == 0, "early completion accrued coins")
	x.finishTime(id)
	x.fail(c, "event.v1.EventService/CancelEvent", map[string]any{"id": id}, "failed_precondition")
	x.fail(a, "participation.v1.ParticipationService/CancelParticipation", map[string]any{
		"id": first,
	}, "failed_precondition")
	x.fail(c, "event.v1.EventService/CompleteEvent", map[string]any{
		"id":         id,
		"request_id": "complete",
	}, "invalid_argument")
	x.fail(c, "event.v1.EventService/CompleteEvent", map[string]any{
		"id":         id,
		"request_id": "complete",
		"results": []any{
			map[string]any{
				"participation_id": first,
				"absent":           false,
			},
		},
	}, "invalid_argument")
	results := []any{map[string]any{
		"participation_id":  first,
		"confirmed_minutes": "180",
	}}
	x.call(c, "event.v1.EventService/CompleteEvent", map[string]any{
		"id":         id,
		"request_id": "complete",
		"results":    results,
	})
	x.call(c, "event.v1.EventService/CompleteEvent", map[string]any{
		"id":         id,
		"request_id": "complete",
		"results":    results,
	})
	x.fail(c, "event.v1.EventService/CompleteEvent", map[string]any{
		"id":         id,
		"request_id": "complete",
		"results": []any{
			map[string]any{
				"participation_id": first,
				"absent":           true,
			},
		},
	}, "already_exists")
	x.fail(c, "event.v1.EventService/CompleteEvent", map[string]any{
		"id":         id,
		"request_id": "other",
		"results":    results,
	}, "failed_precondition")
	profile := object(x.t, x.call(a, "volunteer.v1.VolunteerService/GetVolunteer", nil), "volunteer")
	require(
		x.t,
		integer(x.t, profile, "coin_balance") == 30 && integer(x.t, profile, "confirmed_minutes") == 180,
		"completion totals incorrect",
	)
	closed := object(x.t, x.call(waiting, "participation.v1.ParticipationService/GetParticipation", map[string]any{
		"id": pending,
	}), "participation")
	require(x.t, closed["status"] == "PARTICIPATION_STATUS_CLOSED", "pending application not closed")
	history := x.call(a, "participation.v1.ParticipationService/ListParticipationChanges", map[string]any{
		"participation_id": first,
	})
	require(x.t, len(list(x.t, history, "changes")) == 6, "participation audit lost a transition")
	x.fail(b, "participation.v1.ParticipationService/ListParticipationChanges", map[string]any{
		"participation_id": first,
	}, "permission_denied")
	scoped := x.call(c, "participation.v1.ParticipationService/ListParticipations", map[string]any{
		"scope":    2,
		"event_id": id,
		"status":   4,
	})
	require(x.t, len(list(x.t, scoped, "participations")) == 1, "participation status filter failed")
	mine := x.call(a, "participation.v1.ParticipationService/ListParticipations", map[string]any{"scope": 1})
	require(x.t, len(list(x.t, mine, "participations")) == 1, "volunteer participation scope leaked records")

	finished := object(x.t, x.call(c, "event.v1.EventService/GetEvent", map[string]any{"id": id}), "event")
	fromBoundary := x.call(c, "event.v1.EventService/ListEvents", map[string]any{
		"scope":        2,
		"ends_at_from": finished["ends_at"],
		"ends_at_to":   time.Now().UTC().Format(time.RFC3339Nano),
	})
	require(x.t, len(list(x.t, fromBoundary, "events")) == 1, "ends_at_from excluded its boundary")
	toBoundary := x.call(
		c,
		"event.v1.EventService/ListEvents",
		map[string]any{
			"scope":      2,
			"ends_at_to": finished["ends_at"],
		},
	)
	require(x.t, len(list(x.t, toBoundary, "events")) == 0, "ends_at_to included its boundary")
	beyond := x.call(c, "event.v1.EventService/ListEvents", map[string]any{
		"scope":  2,
		"offset": "9223372036854775807",
	})
	require(
		x.t,
		len(list(x.t, beyond, "events")) == 0 && integer(x.t, beyond, "total_count") == 3,
		"out-of-range offset changed total or returned records",
	)
	x.fail(c, "event.v1.EventService/ListEvents", map[string]any{
		"scope":  2,
		"status": 0,
	}, "invalid_argument")
	x.fail(c, "event.v1.EventService/ListEvents", map[string]any{
		"scope":  2,
		"status": 99,
	}, "invalid_argument")
	x.call(a, "cooperation.v1.CooperationService/TerminateCooperation", map[string]any{"coordinator_id": c.id})
	x.call(a, "event.v1.EventService/GetEvent", map[string]any{"id": id})
	x.fail(a, "event.v1.EventService/ListEvents", map[string]any{
		"scope": 1,
		"search_area": map[string]any{
			"city_id": "missing",
		},
	}, "not_found")
	cancelID := x.event(c, "Отменяемая акция", 2)
	application := text(x.t, x.call(b, "participation.v1.ParticipationService/RequestParticipation", map[string]any{
		"event_id": cancelID,
	})["id"])
	x.call(c, "event.v1.EventService/CancelEvent", map[string]any{"id": cancelID})
	canceled := object(x.t, x.call(b, "participation.v1.ParticipationService/GetParticipation", map[string]any{
		"id": application,
	}), "participation")
	require(
		x.t,
		canceled["cancellation_reason"] == "PARTICIPATION_CANCELLATION_REASON_EVENT_CANCELLED",
		"event cancellation reason missing",
	)
	x.fail(c, "event.v1.EventService/UpdateEvent", map[string]any{"id": cancelID}, "failed_precondition")
	x.fail(c, "event.v1.EventService/UpdateEvent", map[string]any{
		"id":      x.event(c, "Некратная дата", 1),
		"ends_at": time.Now().Add(70 * time.Minute).Format(time.RFC3339Nano),
	}, "invalid_argument")
	absentEvent := x.event(c, "Неявка", 1)
	absentID := x.participate(c, b, absentEvent)
	x.finishTime(absentEvent)
	x.call(c, "event.v1.EventService/CompleteEvent", map[string]any{
		"id":         absentEvent,
		"request_id": "absent",
		"results": []any{map[string]any{
			"participation_id": absentID,
			"absent":           true,
		}},
	})
	absent := object(
		x.t,
		x.call(b, "participation.v1.ParticipationService/GetParticipation", map[string]any{"id": absentID}),
		"participation",
	)
	require(
		x.t,
		absent["status"] == "PARTICIPATION_STATUS_NOT_ATTENDED" && absent["awarded_coins"] == nil,
		"absent volunteer received an award",
	)
	absentProfile := object(x.t, x.call(b, "volunteer.v1.VolunteerService/GetVolunteer", nil), "volunteer")
	require(
		x.t,
		integer(x.t, absentProfile, "coin_balance") == 0 && integer(x.t, absentProfile, "confirmed_minutes") == 0,
		"absence changed wallet",
	)
	x.call(b, "event_chat.v1.EventChatService/GetEventChat", map[string]any{"event_id": absentEvent})
}
