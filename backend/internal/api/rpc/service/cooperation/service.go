package rpc_service_cooperation

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	cooperation_v1 "github.com/Roum1212/hahaton/genproto/cooperation/v1"
	cooperation_v1connect "github.com/Roum1212/hahaton/genproto/cooperation/v1/cooperation_v1connect"
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
)

var _ cooperation_v1connect.CooperationServiceHandler = Service{}

type Service struct {
	acceptCooperationInvitationRPC rpc_service_cooperation_rpc_accept_cooperation_invitation.RPC
	approveCooperationRPC          rpc_service_cooperation_rpc_approve_cooperation.RPC
	createCooperationInvitationRPC rpc_service_cooperation_rpc_create_cooperation_invitation.RPC
	getCooperationRPC              rpc_service_cooperation_rpc_get_cooperation.RPC
	listCooperationInvitationsRPC  rpc_service_cooperation_rpc_list_cooperation_invitations.RPC
	listCooperationsRPC            rpc_service_cooperation_rpc_list_cooperations.RPC
	rejectCooperationRPC           rpc_service_cooperation_rpc_reject_cooperation.RPC
	requestCooperationRPC          rpc_service_cooperation_rpc_request_cooperation.RPC
	revokeCooperationInvitationRPC rpc_service_cooperation_rpc_revoke_cooperation_invitation.RPC
	terminateCooperationRPC        rpc_service_cooperation_rpc_terminate_cooperation.RPC
}

func (x Service) AcceptCooperationInvitation(
	ctx context.Context,
	r *connect.Request[cooperation_v1.AcceptCooperationInvitationRequest],
) (*connect.Response[cooperation_v1.AcceptCooperationInvitationResponse], error) {
	return x.acceptCooperationInvitationRPC.AcceptCooperationInvitation(ctx, r)
}

func (x Service) ApproveCooperation(
	ctx context.Context,
	r *connect.Request[cooperation_v1.ApproveCooperationRequest],
) (*connect.Response[emptypb.Empty], error) {
	return x.approveCooperationRPC.ApproveCooperation(ctx, r)
}

func (x Service) CreateCooperationInvitation(
	ctx context.Context,
	r *connect.Request[cooperation_v1.CreateCooperationInvitationRequest],
) (*connect.Response[cooperation_v1.CreateCooperationInvitationResponse], error) {
	return x.createCooperationInvitationRPC.CreateCooperationInvitation(ctx, r)
}

func (x Service) GetCooperation(
	ctx context.Context,
	r *connect.Request[cooperation_v1.GetCooperationRequest],
) (*connect.Response[cooperation_v1.GetCooperationResponse], error) {
	return x.getCooperationRPC.GetCooperation(ctx, r)
}

func (x Service) ListCooperationInvitations(
	ctx context.Context,
	r *connect.Request[cooperation_v1.ListCooperationInvitationsRequest],
) (*connect.Response[cooperation_v1.ListCooperationInvitationsResponse], error) {
	return x.listCooperationInvitationsRPC.ListCooperationInvitations(ctx, r)
}

func (x Service) ListCooperations(
	ctx context.Context,
	r *connect.Request[cooperation_v1.ListCooperationsRequest],
) (*connect.Response[cooperation_v1.ListCooperationsResponse], error) {
	return x.listCooperationsRPC.ListCooperations(ctx, r)
}

func (x Service) RejectCooperation(
	ctx context.Context,
	r *connect.Request[cooperation_v1.RejectCooperationRequest],
) (*connect.Response[emptypb.Empty], error) {
	return x.rejectCooperationRPC.RejectCooperation(ctx, r)
}

func (x Service) RequestCooperation(
	ctx context.Context,
	r *connect.Request[cooperation_v1.RequestCooperationRequest],
) (*connect.Response[cooperation_v1.RequestCooperationResponse], error) {
	return x.requestCooperationRPC.RequestCooperation(ctx, r)
}

func (x Service) RevokeCooperationInvitation(
	ctx context.Context,
	r *connect.Request[cooperation_v1.RevokeCooperationInvitationRequest],
) (*connect.Response[emptypb.Empty], error) {
	return x.revokeCooperationInvitationRPC.RevokeCooperationInvitation(ctx, r)
}

func (x Service) TerminateCooperation(
	ctx context.Context,
	r *connect.Request[cooperation_v1.TerminateCooperationRequest],
) (*connect.Response[emptypb.Empty], error) {
	return x.terminateCooperationRPC.TerminateCooperation(ctx, r)
}

func NewService(
	acceptCooperationInvitationRPC rpc_service_cooperation_rpc_accept_cooperation_invitation.RPC,
	approveCooperationRPC rpc_service_cooperation_rpc_approve_cooperation.RPC,
	createCooperationInvitationRPC rpc_service_cooperation_rpc_create_cooperation_invitation.RPC,
	getCooperationRPC rpc_service_cooperation_rpc_get_cooperation.RPC,
	listCooperationInvitationsRPC rpc_service_cooperation_rpc_list_cooperation_invitations.RPC,
	listCooperationsRPC rpc_service_cooperation_rpc_list_cooperations.RPC,
	rejectCooperationRPC rpc_service_cooperation_rpc_reject_cooperation.RPC,
	requestCooperationRPC rpc_service_cooperation_rpc_request_cooperation.RPC,
	revokeCooperationInvitationRPC rpc_service_cooperation_rpc_revoke_cooperation_invitation.RPC,
	terminateCooperationRPC rpc_service_cooperation_rpc_terminate_cooperation.RPC,
) Service {
	return Service{
		acceptCooperationInvitationRPC: acceptCooperationInvitationRPC,
		approveCooperationRPC:          approveCooperationRPC,
		createCooperationInvitationRPC: createCooperationInvitationRPC,
		getCooperationRPC:              getCooperationRPC,
		listCooperationInvitationsRPC:  listCooperationInvitationsRPC,
		listCooperationsRPC:            listCooperationsRPC,
		rejectCooperationRPC:           rejectCooperationRPC,
		requestCooperationRPC:          requestCooperationRPC,
		revokeCooperationInvitationRPC: revokeCooperationInvitationRPC,
		terminateCooperationRPC:        terminateCooperationRPC,
	}
}
