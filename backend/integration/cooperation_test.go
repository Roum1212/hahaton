//go:build integration

package integration

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func checkCooperation(x *fixture) {
	c, d := x.coordinator("coordinator"), x.coordinator("other-coordinator")
	a, b := x.login("volunteer"), x.login("other-volunteer")
	x.fail(c, "cooperation.v1.CooperationService/RequestCooperation", map[string]any{
		"coordinator_id": c.id,
	}, "invalid_argument")
	x.fail(a, "cooperation.v1.CooperationService/ListCooperations", map[string]any{"scope": 2}, "permission_denied")
	x.call(a, "cooperation.v1.CooperationService/RequestCooperation", map[string]any{"coordinator_id": c.id})
	x.fail(a, "cooperation.v1.CooperationService/RequestCooperation", map[string]any{
		"coordinator_id": c.id,
	}, "already_exists")
	x.call(c, "cooperation.v1.CooperationService/RejectCooperation", map[string]any{"volunteer_id": a.id})
	rejected := object(x.t, x.call(a, "cooperation.v1.CooperationService/GetCooperation", map[string]any{
		"coordinator_id": c.id,
	}), "cooperation")
	retry := rejected["retry_after"]

	x.call(c, "coordinator.v1.CoordinatorService/UpdateCoordinatorSettings", map[string]any{
		"cooperation_request_cooldown_seconds": "0",
	})
	unchanged := object(x.t, x.call(c, "cooperation.v1.CooperationService/GetCooperation", map[string]any{
		"volunteer_id": a.id,
	}), "cooperation")
	require(x.t, unchanged["retry_after"] == retry, "cooldown update changed existing retry_after")
	x.fail(a, "cooperation.v1.CooperationService/RequestCooperation", map[string]any{
		"coordinator_id": c.id,
	}, "failed_precondition")
	invitation := object(
		x.t,
		x.call(c, "cooperation.v1.CooperationService/CreateCooperationInvitation", nil),
		"invitation",
	)
	id := text(x.t, invitation["id"])
	url := text(x.t, invitation["invitation_url"])
	_, token, ok := strings.Cut(url, "/invite/")
	require(x.t, ok, "invitation URL does not contain token")

	var stored string
	must(x.t, x.pool.QueryRow(x.t.Context(), "SELECT token FROM cooperation_invitations WHERE id=$1", id).Scan(&stored))
	require(x.t, stored != token && !strings.Contains(stored, token), "invitation secret stored in plaintext")
	listed := x.call(c, "cooperation.v1.CooperationService/ListCooperationInvitations", nil)
	values := list(x.t, listed, "invitations")
	require(x.t, len(values) == 1, "invitation missing from list")
	require(x.t, dictionary(x.t, values[0])["invitation_url"] == url, "invitation URL changed between requests")
	x.fail(d, "cooperation.v1.CooperationService/RevokeCooperationInvitation", map[string]any{
		"id": id,
	}, "permission_denied")
	x.fail(c, "cooperation.v1.CooperationService/AcceptCooperationInvitation", map[string]any{
		"token": token,
	}, "invalid_argument")
	accepted := object(x.t, x.call(a, "cooperation.v1.CooperationService/AcceptCooperationInvitation", map[string]any{
		"token": token,
	}), "cooperation")
	repeated := object(x.t, x.call(a, "cooperation.v1.CooperationService/AcceptCooperationInvitation", map[string]any{
		"token": token,
	}), "cooperation")
	require(x.t, accepted["updated_at"] == repeated["updated_at"], "active acceptance changed status timestamp")
	x.call(b, "cooperation.v1.CooperationService/AcceptCooperationInvitation", map[string]any{"token": token})
	x.call(c, "cooperation.v1.CooperationService/RevokeCooperationInvitation", map[string]any{"id": id})
	x.call(c, "cooperation.v1.CooperationService/RevokeCooperationInvitation", map[string]any{"id": id})
	missing := x.fail(a, "cooperation.v1.CooperationService/AcceptCooperationInvitation", map[string]any{
		"token": token,
	}, "not_found")
	require(x.t, resourceInfo(x.t, missing).GetResourceName() == "", "invitation token exposed in error")
	x.call(c, "cooperation.v1.CooperationService/TerminateCooperation", map[string]any{"volunteer_id": b.id})
	x.call(b, "cooperation.v1.CooperationService/RequestCooperation", map[string]any{"coordinator_id": c.id})
	x.call(c, "cooperation.v1.CooperationService/ApproveCooperation", map[string]any{"volunteer_id": b.id})
	x.fail(d, "volunteer.v1.VolunteerService/GetCoordinatorVolunteer", map[string]any{
		"volunteer_id": a.id,
	}, "permission_denied")
	contacts := object(x.t, x.call(c, "volunteer.v1.VolunteerService/GetCoordinatorVolunteer", map[string]any{
		"volunteer_id": a.id,
	}), "volunteer")
	require(x.t, contacts["email_address"] == "volunteer@example.com", "active cooperation contacts missing")
	require(x.t, contacts["coin_balance"] == nil && contacts["date_of_birth"] == nil, "private volunteer data leaked")
	x.call(a, "cooperation.v1.CooperationService/TerminateCooperation", map[string]any{"coordinator_id": c.id})
	contacts = object(x.t, x.call(c, "volunteer.v1.VolunteerService/GetCoordinatorVolunteer", map[string]any{
		"volunteer_id": a.id,
	}), "volunteer")
	require(x.t, contacts["email_address"] == nil &&
		contacts["phone_number"] == nil, "contacts survived cooperation termination")
	listResponse := x.call(c, "cooperation.v1.CooperationService/ListCooperations", map[string]any{
		"scope":  2,
		"status": 4,
		"limit":  1,
	})
	require(x.t, fmt.Sprint(listResponse["total_count"]) == "1", "cooperation status filter incorrect")
	x.fail(a, "cooperation.v1.CooperationService/GetCooperation", nil, "invalid_argument")
	x.fail(c, "cooperation.v1.CooperationService/CreateCooperationInvitation", map[string]any{
		"expires_at": time.Now().Add(-time.Second).Format(time.RFC3339Nano),
	}, "invalid_argument")
	denied := x.request(actor{}, "admin.v1.AdminService/FindVolunteers", nil, true)
	require(x.t, denied.Code == http.StatusBadRequest, "empty admin search criterion accepted")

	var changes int
	must(x.t, x.pool.QueryRow(
		x.t.Context(),
		"SELECT count(*) FROM cooperation_changes WHERE coordinator_id=$1 AND volunteer_id=$2",
		c.id,
		a.id,
	).Scan(
		&changes,
	))
	require(x.t, changes == 4, "cooperation audit lost transitions or duplicated active acceptance")
}
