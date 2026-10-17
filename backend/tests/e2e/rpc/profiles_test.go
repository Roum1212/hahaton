//go:build e2e

package rpc

import (
	"crypto/rand"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	admin_v1 "github.com/Roum1212/hahaton/genproto/admin/v1"
	coordinator_v1 "github.com/Roum1212/hahaton/genproto/coordinator/v1"
	volunteer_v1 "github.com/Roum1212/hahaton/genproto/volunteer/v1"
)

func (x *RPCTestSuite) TestAdminProvisioning() {
	ctx := x.T().Context()
	subject := strings.ToLower(rand.Text())
	a := x.loginSubject(subject, nil)
	r := connect.NewRequest(
		&admin_v1.CreateCoordinatorRequest{
			VolunteerId:                       a.id,
			CooperationRequestCooldownSeconds: proto.Int64(60),
		},
	)
	_, err := x.clients.admin.CreateCoordinator(ctx, r)
	callError(x.T(), err, connect.CodeUnauthenticated)
	r.Header().Set("Authorization", "Bearer incorrect")
	_, err = x.clients.admin.CreateCoordinator(ctx, r)
	callError(x.T(), err, connect.CodeUnauthenticated)
	r.Header().Set("Authorization", "Bearer "+adminAPIKey)
	created, err := x.clients.admin.CreateCoordinator(ctx, r)
	x.Require().NoError(err)
	x.Require().Equal(a.id, created.Msg.GetId())
	_, err = x.clients.admin.CreateCoordinator(ctx, r)
	callError(x.T(), err, connect.CodeAlreadyExists)
	current, err := x.clients.session.GetSession(ctx, request(a, &emptypb.Empty{}))
	x.Require().NoError(err)
	x.Require().True(current.Msg.GetSession().GetHasCoordinatorProfile())

	cases := []*admin_v1.FindVolunteersRequest{
		{Criterion: &admin_v1.FindVolunteersRequest_VolunteerId{VolunteerId: a.id}},
		{
			Criterion: &admin_v1.FindVolunteersRequest_EmailAddress{
				EmailAddress: " " + strings.ToUpper(subject) + "@EXAMPLE.COM ",
			},
		},
		{Criterion: &admin_v1.FindVolunteersRequest_PhoneNumber{PhoneNumber: "+79990000001"}},
	}
	for _, message := range cases {
		lookup := connect.NewRequest(message)
		lookup.Header().Set("Authorization", "Bearer "+adminAPIKey)
		found, err := x.clients.admin.FindVolunteers(ctx, lookup)
		x.Require().NoError(err)

		ids := make([]string, 0, len(found.Msg.GetVolunteers()))
		for _, volunteer := range found.Msg.GetVolunteers() {
			ids = append(ids, volunteer.GetId())
		}

		x.Require().Contains(ids, a.id)
	}

	empty := connect.NewRequest(&admin_v1.FindVolunteersRequest{})
	empty.Header().Set("Authorization", "Bearer "+adminAPIKey)
	_, err = x.clients.admin.FindVolunteers(ctx, empty)
	callError(x.T(), err, connect.CodeInvalidArgument)
	publicAdmin := x.baseURL + "/admin.v1.AdminService/FindVolunteers"
	publicRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, publicAdmin, strings.NewReader("{}"))
	x.Require().NoError(err)
	publicRequest.Header.Set("Content-Type", "application/json")
	response, err := x.httpClient.Do(publicRequest)
	x.Require().NoError(err)
	x.Require().Equal(404, response.StatusCode)
	x.Require().NoError(response.Body.Close())
}

func (x *RPCTestSuite) TestCoordinatorProfileAndSettings() {
	ctx := x.T().Context()
	c, outsider := x.coordinator(), x.login()
	_, err := x.clients.coordinator.GetCoordinatorSettings(ctx,
		request(outsider, &coordinator_v1.GetCoordinatorSettingsRequest{}))
	callError(x.T(), err, connect.CodePermissionDenied)
	_, err = x.clients.coordinator.UpdateCoordinator(ctx, request(outsider, &coordinator_v1.UpdateCoordinatorRequest{
		Description: proto.String("forbidden"),
	}))
	callError(x.T(), err, connect.CodePermissionDenied)
	_, err = x.clients.coordinator.UpdateCoordinator(ctx, request(c, &coordinator_v1.UpdateCoordinatorRequest{
		Description: proto.String("Описание"),
	}))
	x.Require().NoError(err)
	_, err = x.clients.coordinator.UpdateCoordinator(ctx, request(c, &coordinator_v1.UpdateCoordinatorRequest{}))
	x.Require().NoError(err)
	profile, err := x.clients.coordinator.GetCoordinator(ctx, request(outsider, &coordinator_v1.GetCoordinatorRequest{
		Id: proto.String(c.id),
	}))
	x.Require().NoError(err)
	x.Require().Equal("Описание", profile.Msg.GetCoordinator().GetDescription())
	_, err = x.clients.volunteer.UpdateVolunteer(ctx, request(c, &volunteer_v1.UpdateVolunteerRequest{
		FirstName: proto.String("Новое имя"),
	}))
	x.Require().NoError(err)
	updated, err := x.clients.coordinator.GetCoordinator(ctx, request(c, &coordinator_v1.GetCoordinatorRequest{}))
	x.Require().NoError(err)
	x.Require().Contains(updated.Msg.GetCoordinator().GetDisplayName(), "Новое имя")
	_, err = x.clients.coordinator.UpdateCoordinator(ctx, request(c, &coordinator_v1.UpdateCoordinatorRequest{
		Description: proto.String(""),
	}))
	x.Require().NoError(err)
	cleared, err := x.clients.coordinator.GetCoordinator(ctx, request(c, &coordinator_v1.GetCoordinatorRequest{}))
	x.Require().NoError(err)
	x.Require().Empty(cleared.Msg.GetCoordinator().GetDescription())

	for _, seconds := range []int64{
		60,
		0,
	} {
		_, err = x.clients.coordinator.UpdateCoordinatorSettings(ctx, request(c,
			&coordinator_v1.UpdateCoordinatorSettingsRequest{CooperationRequestCooldownSeconds: proto.Int64(seconds)}))
		x.Require().NoError(err)
		_, err = x.clients.coordinator.UpdateCoordinatorSettings(ctx, request(c,
			&coordinator_v1.UpdateCoordinatorSettingsRequest{}))
		x.Require().NoError(err)
		settings, err := x.clients.coordinator.GetCoordinatorSettings(ctx,
			request(c, &coordinator_v1.GetCoordinatorSettingsRequest{}))
		x.Require().NoError(err)
		x.Require().Equal(seconds, settings.Msg.GetSettings().GetCooperationRequestCooldownSeconds())
	}

	_, err = x.clients.coordinator.UpdateCoordinatorSettings(ctx, request(c,
		&coordinator_v1.UpdateCoordinatorSettingsRequest{CooperationRequestCooldownSeconds: proto.Int64(-1)}))
	callError(x.T(), err, connect.CodeInvalidArgument)
	listed, err := x.clients.coordinator.ListCoordinators(ctx,
		request(outsider, &coordinator_v1.ListCoordinatorsRequest{Limit: 1}))
	x.Require().NoError(err)
	x.Require().Len(listed.Msg.GetCoordinators(), 1)
	x.Require().GreaterOrEqual(listed.Msg.GetTotalCount(), int64(1))
}

func (x *RPCTestSuite) TestVolunteerProfileAndRegistration() {
	ctx := x.T().Context()
	subject := strings.ToLower(rand.Text())
	a := x.loginSubject(subject, map[string]any{
		"email":         "John.Doe@googlemail.com",
		"date_of_birth": "2000-02-29",
	})
	b := x.login()
	initial, err := x.clients.volunteer.GetVolunteer(ctx, request(a, &volunteer_v1.GetVolunteerRequest{}))
	x.Require().NoError(err)
	x.Require().Equal("johndoe@gmail.com", initial.Msg.GetVolunteer().GetEmailAddress())
	x.Require().EqualValues(29, initial.Msg.GetVolunteer().GetDateOfBirth().GetDay())
	x.Require().Zero(initial.Msg.GetVolunteer().GetCoinBalance())
	_, err = x.clients.volunteer.GetVolunteer(
		ctx,
		request(b, &volunteer_v1.GetVolunteerRequest{Id: proto.String(a.id)}),
	)
	callError(x.T(), err, connect.CodePermissionDenied)
	_, err = x.clients.volunteer.UpdateVolunteer(ctx, request(a, &volunteer_v1.UpdateVolunteerRequest{
		FirstName:  proto.String("Пётр"),
		MiddleName: proto.String(""),
	}))
	x.Require().NoError(err)
	_, err = x.clients.volunteer.UpdateVolunteer(ctx, request(a, &volunteer_v1.UpdateVolunteerRequest{}))
	x.Require().NoError(err)
	again := x.loginSubject(subject, map[string]any{
		"given_name":    nil,
		"email":         nil,
		"date_of_birth": "invalid",
	})
	x.Require().Equal(a.id, again.id)
	updated, err := x.clients.volunteer.GetVolunteer(ctx, request(again, &volunteer_v1.GetVolunteerRequest{}))
	x.Require().NoError(err)
	x.Require().Equal("Пётр", updated.Msg.GetVolunteer().GetFirstName())
	x.Require().Empty(updated.Msg.GetVolunteer().GetMiddleName())
	x.Require().Equal("Иванов", updated.Msg.GetVolunteer().GetLastName())
	_, err = x.clients.volunteer.UpdateVolunteer(ctx, request(a, &volunteer_v1.UpdateVolunteerRequest{
		DateOfBirth: &date.Date{
			Year:  2001,
			Month: 2,
			Day:   29,
		},
	}))
	callError(x.T(), err, connect.CodeInvalidArgument)

	var rpcError *connect.Error
	x.Require().ErrorAs(err, &rpcError)
	x.Require().NotEmpty(rpcError.Details())
	detail, err := rpcError.Details()[0].Value()
	x.Require().NoError(err)

	badRequest, ok := detail.(*errdetails.BadRequest)
	x.Require().True(ok)
	x.Require().NotEmpty(badRequest.GetFieldViolations()[0].GetField())
	history, err := x.clients.volunteer.ListVolunteerOperations(ctx,
		request(a, &volunteer_v1.ListVolunteerOperationsRequest{}))
	x.Require().NoError(err)
	x.Require().Empty(history.Msg.GetOperations())
	x.Require().Zero(history.Msg.GetTotalCount())
}
