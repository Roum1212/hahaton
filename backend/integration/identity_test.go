//go:build integration

package integration

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func checkIdentity(x *fixture) {
	a, b := x.login("alice"), x.login("bob")
	require(x.t, a.id != b.id, "different subjects merged")
	x.call(a, "session.v1.SessionService/OpenSession", nil)
	current := object(x.t, x.call(a, "session.v1.SessionService/GetSession", nil), "session")
	require(x.t, current["user_id"] == a.id, "open logged out an authenticated user")
	require(x.t, current["has_coordinator_profile"] != true, "volunteer acquired coordinator role")
	x.fail(actor{}, "volunteer.v1.VolunteerService/GetVolunteer", nil, "unauthenticated")
	opened := x.request(actor{}, "session.v1.SessionService/OpenSession", nil, false)
	require(x.t, opened.Code == http.StatusOK, "anonymous session was not opened")
	anonymous := actor{cookie: opened.Result().Cookies()[0]}
	x.fail(anonymous, "volunteer.v1.VolunteerService/GetVolunteer", nil, "unauthenticated")
	x.fail(anonymous, "event.v1.EventService/CreateEvent", nil, "unauthenticated")
	x.fail(a, "volunteer.v1.VolunteerService/GetVolunteer", map[string]any{"id": b.id}, "permission_denied")
	details := x.fail(a, "volunteer.v1.VolunteerService/GetVolunteer", map[string]any{"id": ""}, "invalid_argument")
	require(x.t, len(list(x.t, details, "details")) > 0, "invalid argument details missing")
	missing := x.fail(a, "coordinator.v1.CoordinatorService/GetCoordinator", nil, "not_found")
	info := resourceInfo(x.t, missing)
	require(x.t, info.GetResourceType() == "coordinator.v1.Coordinator", "wrong missing coordinator resource")
	require(x.t, info.GetResourceName() == "id:"+a.id, "missing coordinator id not included")
	x.fail(a, "coordinator.v1.CoordinatorService/GetCoordinatorSettings", nil, "permission_denied")
	x.fail(a, "coordinator.v1.CoordinatorService/UpdateCoordinatorSettings", nil, "permission_denied")
	x.fail(a, "coordinator.v1.CoordinatorService/UpdateCoordinator", map[string]any{
		"description": "test",
	}, "permission_denied")
	x.call(a, "volunteer.v1.VolunteerService/UpdateVolunteer", map[string]any{
		"first_name":  "Новое имя",
		"middle_name": "Отчество",
		"date_of_birth": map[string]any{
			"year":  2000,
			"month": 2,
			"day":   29,
		},
	})
	again := x.login("alice")
	require(x.t, again.id == a.id, "repeat login created another volunteer")
	profile := object(
		x.t,
		x.call(again, "volunteer.v1.VolunteerService/GetVolunteer", map[string]any{"id": a.id}),
		"volunteer",
	)
	require(x.t, profile["first_name"] == "Новое имя", "repeat login overwrote profile")
	require(x.t, integer(x.t, profile, "coin_balance") == 0, "new volunteer has nonzero balance")
	x.fail(a, "volunteer.v1.VolunteerService/UpdateVolunteer", map[string]any{"date_of_birth": map[string]any{
		"year":  2001,
		"month": 2,
		"day":   29,
	}}, "invalid_argument")
	x.fail(a, "volunteer.v1.VolunteerService/UpdateVolunteer", map[string]any{"first_name": " "}, "invalid_argument")

	found := x.request(actor{}, "admin.v1.AdminService/FindVolunteers", map[string]any{
		"email_address": " ALICE@EXAMPLE.COM ",
	}, true)
	require(x.t, found.Code == http.StatusOK, found.Body.String())
	foundValues := list(x.t, document(x.t, found.Body.Bytes()), "volunteers")
	require(x.t, len(foundValues) == 1, "email lookup not normalized")

	for _, criterion := range []map[string]any{
		{"volunteer_id": a.id},
		{"phone_number": "+79990000001"},
	} {
		response := x.request(actor{}, "admin.v1.AdminService/FindVolunteers", criterion, true)
		require(x.t, response.Code == http.StatusOK, response.Body.String())
	}

	denied := x.request(a, "admin.v1.AdminService/FindVolunteers", map[string]any{"volunteer_id": a.id}, false)
	require(x.t, denied.Code == http.StatusNotFound, "admin route exposed on client listener")
	request := x.request(actor{}, "admin.v1.AdminService/CreateCoordinator", map[string]any{
		"volunteer_id":                         a.id,
		"cooperation_request_cooldown_seconds": "60",
	}, true)
	require(x.t, request.Code == http.StatusOK, request.Body.String())
	duplicate := x.request(actor{}, "admin.v1.AdminService/CreateCoordinator", map[string]any{
		"volunteer_id":                         a.id,
		"cooperation_request_cooldown_seconds": "60",
	}, true)
	require(x.t, duplicate.Code == http.StatusConflict, "duplicate coordinator accepted")
	current = object(x.t, x.call(a, "session.v1.SessionService/GetSession", nil), "session")
	require(x.t, current["has_coordinator_profile"] == true, "coordinator rights not visible in current session")
	coordinator := object(
		x.t,
		x.call(b, "coordinator.v1.CoordinatorService/GetCoordinator", map[string]any{"id": a.id}),
		"coordinator",
	)
	require(
		x.t,
		strings.Contains(text(x.t, coordinator["display_name"]), "Новое имя"),
		"coordinator name is not shared",
	)
	x.call(a, "coordinator.v1.CoordinatorService/UpdateCoordinator", map[string]any{"description": "Описание"})
	x.call(a, "coordinator.v1.CoordinatorService/UpdateCoordinatorSettings", map[string]any{
		"cooperation_request_cooldown_seconds": "0",
	})
	x.call(a, "coordinator.v1.CoordinatorService/UpdateCoordinatorSettings", nil)
	settings := object(x.t, x.call(a, "coordinator.v1.CoordinatorService/GetCoordinatorSettings", nil), "settings")
	require(
		x.t,
		integer(x.t, settings, "cooperation_request_cooldown_seconds") == 0,
		"settings absent update changed zero",
	)
	x.fail(a, "coordinator.v1.CoordinatorService/UpdateCoordinatorSettings", map[string]any{
		"cooperation_request_cooldown_seconds": "-1",
	}, "invalid_argument")
	coordinators := x.call(b, "coordinator.v1.CoordinatorService/ListCoordinators", nil)
	require(x.t, fmt.Sprint(coordinators["total_count"]) == "1", "coordinator catalog incorrect")
	checkListInputs(x, a)

	x.call(a, "session.v1.SessionService/CloseSession", nil)
	x.fail(a, "session.v1.SessionService/GetSession", nil, "unauthenticated")
	x.call(actor{}, "session.v1.SessionService/CloseSession", nil)
	checkInvalidTokens(x)
	checkProfileModels(x)
}

func checkInvalidTokens(x *fixture) {
	opened := x.request(actor{}, "session.v1.SessionService/OpenSession", nil, false)
	anonymous := actor{cookie: opened.Result().Cookies()[0]}
	base := map[string]any{
		"iss": x.provider.URL,
		"aud": "demo",
		"sub": "claims",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	variants := []map[string]any{
		{"aud": "other"},
		{"iss": "https://invalid.example.test"},
		{"exp": time.Now().Add(-time.Hour).Unix()},
		{"iat": time.Now().Add(time.Hour).Unix()},
		{"nbf": time.Now().Add(time.Hour).Unix()},
		{"aud": []string{
			"demo",
			"other",
		}},
		{"azp": "other"},
		{"sub": ""},
	}
	for _, changes := range variants {
		claims := make(map[string]any, len(base))
		for key, value := range base {
			claims[key] = value
		}

		for key, value := range changes {
			claims[key] = value
		}

		x.fail(
			anonymous,
			"session.v1.SessionService/AuthenticateSession",
			map[string]any{"id_token": x.token(claims)},
			"unauthenticated",
		)
	}

	x.fail(anonymous, "session.v1.SessionService/AuthenticateSession", map[string]any{
		"id_token": "invalid",
	}, "unauthenticated")
	x.fail(
		anonymous,
		"session.v1.SessionService/AuthenticateSession",
		map[string]any{"id_token": x.token(base)},
		"failed_precondition",
	)
}

func checkProfileModels(x *fixture) {
	cases := []struct {
		email     string
		canonical string
	}{
		{
			"John.Doe@googlemail.com",
			"johndoe@gmail.com",
		},
		{
			"John-Doe@ya.ru",
			"john.doe@yandex.ru",
		},
		{
			"John_Doe@protonmail.com",
			"johndoe@proton.me",
		},
		{
			"John@me.com",
			"john@icloud.com",
		},
	}
	for index, scenario := range cases {
		subject := fmt.Sprintf("canonical-%d", index)
		a := x.login(subject, map[string]any{
			"email":         scenario.email,
			"date_of_birth": "2000-02-29",
		})
		profile := object(x.t, x.call(a, "volunteer.v1.VolunteerService/GetVolunteer", nil), "volunteer")
		require(x.t, profile["email_address"] == scenario.canonical, "flida email normalization differs")
		date := object(x.t, profile, "date_of_birth")
		require(x.t, integer(x.t, date, "day") == 29, "leap date was changed in storage")

		again := x.login(subject, map[string]any{
			"given_name":    nil,
			"email":         nil,
			"phone_number":  nil,
			"date_of_birth": "invalid",
		})
		require(x.t, again.id == a.id, "incomplete claims prevented repeat login")
	}

	opened := x.request(actor{}, "session.v1.SessionService/OpenSession", nil, false)
	anonymous := actor{cookie: opened.Result().Cookies()[0]}

	invalid := []map[string]any{
		{"date_of_birth": ""},
		{"date_of_birth": "2001-02-29"},
		{"date_of_birth": "1899-01-01"},
		{"date_of_birth": time.Now().AddDate(1, 0, 0).Format(time.DateOnly)},
		{"email": "invalid@example.test"},
		{"phone_number": "79990000001"},
	}
	for _, overrides := range invalid {
		claims := map[string]any{
			"iss":           x.provider.URL,
			"aud":           "demo",
			"sub":           "invalid-profile",
			"iat":           time.Now().Unix(),
			"exp":           time.Now().Add(time.Hour).Unix(),
			"given_name":    "Имя",
			"family_name":   "Фамилия",
			"email":         "valid@example.com",
			"phone_number":  "+79990000001",
			"date_of_birth": "2000-01-01",
		}
		for key, value := range overrides {
			claims[key] = value
		}

		x.fail(anonymous, "session.v1.SessionService/AuthenticateSession", map[string]any{
			"id_token": x.token(claims),
		}, "failed_precondition")
	}
}

func checkListInputs(x *fixture, a actor) {
	cases := []struct {
		method string
		body   map[string]any
	}{
		{
			"coordinator.v1.CoordinatorService/ListCoordinators",
			nil,
		},
		{
			"cooperation.v1.CooperationService/ListCooperations",
			map[string]any{"scope": 1},
		},
		{
			"event.v1.EventService/ListEvents",
			map[string]any{"scope": 2},
		},
		{
			"participation.v1.ParticipationService/ListParticipations",
			map[string]any{"scope": 1},
		},
		{
			"reward.v1.RewardService/ListRewards",
			map[string]any{"scope": 2},
		},
		{
			"reward_order.v1.RewardOrderService/ListRewardOrders",
			map[string]any{"scope": 1},
		},
		{
			"volunteer.v1.VolunteerService/ListVolunteerOperations",
			nil,
		},
		{
			"cooperation.v1.CooperationService/ListCooperationInvitations",
			nil,
		},
		{
			"event_chat.v1.EventChatService/ListMyEventChats",
			nil,
		},
	}
	for _, scenario := range cases {
		for key, value := range map[string]any{
			"limit":  101,
			"offset": "-1",
		} {
			body := make(map[string]any)
			for name, item := range scenario.body {
				body[name] = item
			}

			body[key] = value
			x.fail(a, scenario.method, body, "invalid_argument")
		}
	}
}
