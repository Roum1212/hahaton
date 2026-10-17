//go:build e2e

package rpc

import (
	"strings"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	cooperation_v1 "github.com/Roum1212/hahaton/genproto/cooperation/v1"
	coordinator_v1 "github.com/Roum1212/hahaton/genproto/coordinator/v1"
	volunteer_v1 "github.com/Roum1212/hahaton/genproto/volunteer/v1"
)

func (x *RPCTestSuite) TestCooperationLifecycleAndPrivacy() {
	ctx := x.T().Context()
	c, other, a := x.coordinator(), x.coordinator(), x.login()
	_, err := x.clients.cooperation.RequestCooperation(ctx, request(c,
		&cooperation_v1.RequestCooperationRequest{CoordinatorId: c.id}))
	callError(x.T(), err, connect.CodeInvalidArgument)
	_, err = x.clients.coordinator.UpdateCoordinatorSettings(ctx, request(c,
		&coordinator_v1.UpdateCoordinatorSettingsRequest{CooperationRequestCooldownSeconds: proto.Int64(60)}))
	x.Require().NoError(err)
	_, err = x.clients.cooperation.RequestCooperation(ctx, request(a,
		&cooperation_v1.RequestCooperationRequest{CoordinatorId: c.id}))
	x.Require().NoError(err)
	_, err = x.clients.cooperation.RequestCooperation(ctx, request(a,
		&cooperation_v1.RequestCooperationRequest{CoordinatorId: c.id}))
	callError(x.T(), err, connect.CodeAlreadyExists)
	_, err = x.clients.cooperation.RejectCooperation(ctx,
		request(c, &cooperation_v1.RejectCooperationRequest{VolunteerId: a.id}))
	x.Require().NoError(err)
	rejected, err := x.clients.cooperation.GetCooperation(ctx, request(a, &cooperation_v1.GetCooperationRequest{
		Counterpart: &cooperation_v1.GetCooperationRequest_CoordinatorId{CoordinatorId: c.id},
	}))
	x.Require().NoError(err)
	x.Require().Equal(
		cooperation_v1.CooperationStatus_COOPERATION_STATUS_REJECTED,
		rejected.Msg.GetCooperation().GetStatus(),
	)
	retry := rejected.Msg.GetCooperation().GetRetryAfter()
	x.Require().True(retry.AsTime().After(time.Now()))
	_, err = x.clients.coordinator.UpdateCoordinatorSettings(ctx, request(c,
		&coordinator_v1.UpdateCoordinatorSettingsRequest{CooperationRequestCooldownSeconds: proto.Int64(0)}))
	x.Require().NoError(err)
	unchanged, err := x.clients.cooperation.GetCooperation(ctx, request(c, &cooperation_v1.GetCooperationRequest{
		Counterpart: &cooperation_v1.GetCooperationRequest_VolunteerId{VolunteerId: a.id},
	}))
	x.Require().NoError(err)
	x.Require().True(proto.Equal(retry, unchanged.Msg.GetCooperation().GetRetryAfter()))
	_, err = x.clients.cooperation.RequestCooperation(ctx, request(a,
		&cooperation_v1.RequestCooperationRequest{CoordinatorId: c.id}))
	callError(x.T(), err, connect.CodeFailedPrecondition)
	invitation, err := x.clients.cooperation.CreateCooperationInvitation(ctx,
		request(c, &cooperation_v1.CreateCooperationInvitationRequest{}))
	x.Require().NoError(err)

	_, token, ok := strings.Cut(invitation.Msg.GetInvitation().GetInvitationUrl(), "/invite/")
	x.Require().True(ok)
	accepted, err := x.clients.cooperation.AcceptCooperationInvitation(ctx,
		request(a, &cooperation_v1.AcceptCooperationInvitationRequest{Token: token}))
	x.Require().NoError(err)
	x.Require().Equal(
		cooperation_v1.CooperationStatus_COOPERATION_STATUS_ACTIVE,
		accepted.Msg.GetCooperation().GetStatus(),
	)
	profile, err := x.clients.volunteer.GetCoordinatorVolunteer(ctx,
		request(c, &volunteer_v1.GetCoordinatorVolunteerRequest{VolunteerId: a.id}))
	x.Require().NoError(err)
	x.Require().NotNil(profile.Msg.GetVolunteer().EmailAddress)
	x.Require().NotNil(profile.Msg.GetVolunteer().PhoneNumber)
	_, err = x.clients.volunteer.GetCoordinatorVolunteer(ctx,
		request(other, &volunteer_v1.GetCoordinatorVolunteerRequest{VolunteerId: a.id}))
	callError(x.T(), err, connect.CodePermissionDenied)
	_, err = x.clients.cooperation.TerminateCooperation(ctx, request(a, &cooperation_v1.TerminateCooperationRequest{
		Counterpart: &cooperation_v1.TerminateCooperationRequest_CoordinatorId{CoordinatorId: c.id},
	}))
	x.Require().NoError(err)
	hidden, err := x.clients.volunteer.GetCoordinatorVolunteer(ctx,
		request(c, &volunteer_v1.GetCoordinatorVolunteerRequest{VolunteerId: a.id}))
	x.Require().NoError(err)
	x.Require().Nil(hidden.Msg.GetVolunteer().EmailAddress)
	x.Require().Nil(hidden.Msg.GetVolunteer().PhoneNumber)
	listed, err := x.clients.cooperation.ListCooperations(ctx, request(c, &cooperation_v1.ListCooperationsRequest{
		Scope:  cooperation_v1.CooperationListScope_COOPERATION_LIST_SCOPE_AS_COORDINATOR,
		Status: cooperation_v1.CooperationStatus_COOPERATION_STATUS_TERMINATED.Enum(),
		Limit:  1,
	}))
	x.Require().NoError(err)
	x.Require().EqualValues(1, listed.Msg.GetTotalCount())
	x.Require().Equal(a.id, listed.Msg.GetCooperations()[0].GetVolunteerId())
	x.cooperate(c, a)
}

func (x *RPCTestSuite) TestReusableInvitationsAndRevocation() {
	ctx := x.T().Context()
	c, other, a, b := x.coordinator(), x.coordinator(), x.login(), x.login()
	created, err := x.clients.cooperation.CreateCooperationInvitation(ctx,
		request(c, &cooperation_v1.CreateCooperationInvitationRequest{}))
	x.Require().NoError(err)

	invitation := created.Msg.GetInvitation()
	_, token, ok := strings.Cut(invitation.GetInvitationUrl(), "/invite/")
	x.Require().True(ok)
	listed, err := x.clients.cooperation.ListCooperationInvitations(ctx,
		request(c, &cooperation_v1.ListCooperationInvitationsRequest{}))
	x.Require().NoError(err)
	x.Require().Len(listed.Msg.GetInvitations(), 1)
	x.Require().Equal(invitation.GetInvitationUrl(), listed.Msg.GetInvitations()[0].GetInvitationUrl())
	_, err = x.clients.cooperation.AcceptCooperationInvitation(ctx,
		request(c, &cooperation_v1.AcceptCooperationInvitationRequest{Token: token}))
	callError(x.T(), err, connect.CodeInvalidArgument)
	first, err := x.clients.cooperation.AcceptCooperationInvitation(ctx,
		request(a, &cooperation_v1.AcceptCooperationInvitationRequest{Token: token}))
	x.Require().NoError(err)
	again, err := x.clients.cooperation.AcceptCooperationInvitation(ctx,
		request(a, &cooperation_v1.AcceptCooperationInvitationRequest{Token: token}))
	x.Require().NoError(err)
	x.Require().True(proto.Equal(first.Msg, again.Msg))
	_, err = x.clients.cooperation.AcceptCooperationInvitation(ctx,
		request(b, &cooperation_v1.AcceptCooperationInvitationRequest{Token: token}))
	x.Require().NoError(err)
	_, err = x.clients.cooperation.RevokeCooperationInvitation(ctx,
		request(other, &cooperation_v1.RevokeCooperationInvitationRequest{Id: invitation.GetId()}))
	callError(x.T(), err, connect.CodePermissionDenied)

	for range 2 {
		_, err = x.clients.cooperation.RevokeCooperationInvitation(ctx,
			request(c, &cooperation_v1.RevokeCooperationInvitationRequest{Id: invitation.GetId()}))
		x.Require().NoError(err)
	}

	_, err = x.clients.cooperation.AcceptCooperationInvitation(ctx,
		request(a, &cooperation_v1.AcceptCooperationInvitationRequest{Token: token}))
	callError(x.T(), err, connect.CodeNotFound)
	_, err = x.clients.cooperation.CreateCooperationInvitation(ctx, request(c,
		&cooperation_v1.CreateCooperationInvitationRequest{ExpiresAt: timestamppb.New(time.Now().Add(-time.Minute))}))
	callError(x.T(), err, connect.CodeInvalidArgument)
	active, err := x.clients.cooperation.GetCooperation(ctx, request(b, &cooperation_v1.GetCooperationRequest{
		Counterpart: &cooperation_v1.GetCooperationRequest_CoordinatorId{CoordinatorId: c.id},
	}))
	x.Require().NoError(err)
	x.Require().Equal(
		cooperation_v1.CooperationStatus_COOPERATION_STATUS_ACTIVE,
		active.Msg.GetCooperation().GetStatus(),
	)
}
