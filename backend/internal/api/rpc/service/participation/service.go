package rpc_service_participation

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	participation_v1 "github.com/Roum1212/hahaton/genproto/participation/v1"
	participation_v1connect "github.com/Roum1212/hahaton/genproto/participation/v1/participation_v1connect"
	rpc_service_participation_rpc_approve_participation "github.com/Roum1212/hahaton/internal/api/rpc/service/participation/rpc/approve-participation"
	rpc_service_participation_rpc_cancel_participation "github.com/Roum1212/hahaton/internal/api/rpc/service/participation/rpc/cancel-participation"
	rpc_service_participation_rpc_get_participation "github.com/Roum1212/hahaton/internal/api/rpc/service/participation/rpc/get-participation"
	rpc_service_participation_rpc_list_participation_changes "github.com/Roum1212/hahaton/internal/api/rpc/service/participation/rpc/list-participation-changes"
	rpc_service_participation_rpc_list_participations "github.com/Roum1212/hahaton/internal/api/rpc/service/participation/rpc/list-participations"
	rpc_service_participation_rpc_reject_participation "github.com/Roum1212/hahaton/internal/api/rpc/service/participation/rpc/reject-participation"
	rpc_service_participation_rpc_request_participation "github.com/Roum1212/hahaton/internal/api/rpc/service/participation/rpc/request-participation"
	rpc_service_participation_rpc_restore_participation "github.com/Roum1212/hahaton/internal/api/rpc/service/participation/rpc/restore-participation"
)

var _ participation_v1connect.ParticipationServiceHandler = Service{}

type Service struct {
	approveParticipationRPC     rpc_service_participation_rpc_approve_participation.RPC
	cancelParticipationRPC      rpc_service_participation_rpc_cancel_participation.RPC
	getParticipationRPC         rpc_service_participation_rpc_get_participation.RPC
	listParticipationChangesRPC rpc_service_participation_rpc_list_participation_changes.RPC
	listParticipationsRPC       rpc_service_participation_rpc_list_participations.RPC
	rejectParticipationRPC      rpc_service_participation_rpc_reject_participation.RPC
	requestParticipationRPC     rpc_service_participation_rpc_request_participation.RPC
	restoreParticipationRPC     rpc_service_participation_rpc_restore_participation.RPC
}

func (x Service) ApproveParticipation(
	ctx context.Context,
	r *connect.Request[participation_v1.ApproveParticipationRequest],
) (*connect.Response[emptypb.Empty], error) {
	return x.approveParticipationRPC.ApproveParticipation(ctx, r)
}

func (x Service) CancelParticipation(
	ctx context.Context,
	r *connect.Request[participation_v1.CancelParticipationRequest],
) (*connect.Response[emptypb.Empty], error) {
	return x.cancelParticipationRPC.CancelParticipation(ctx, r)
}

func (x Service) GetParticipation(
	ctx context.Context,
	r *connect.Request[participation_v1.GetParticipationRequest],
) (*connect.Response[participation_v1.GetParticipationResponse], error) {
	return x.getParticipationRPC.GetParticipation(ctx, r)
}

func (x Service) ListParticipationChanges(
	ctx context.Context,
	r *connect.Request[participation_v1.ListParticipationChangesRequest],
) (*connect.Response[participation_v1.ListParticipationChangesResponse], error) {
	return x.listParticipationChangesRPC.ListParticipationChanges(ctx, r)
}

func (x Service) ListParticipations(
	ctx context.Context,
	r *connect.Request[participation_v1.ListParticipationsRequest],
) (*connect.Response[participation_v1.ListParticipationsResponse], error) {
	return x.listParticipationsRPC.ListParticipations(ctx, r)
}

func (x Service) RejectParticipation(
	ctx context.Context,
	r *connect.Request[participation_v1.RejectParticipationRequest],
) (*connect.Response[emptypb.Empty], error) {
	return x.rejectParticipationRPC.RejectParticipation(ctx, r)
}

func (x Service) RequestParticipation(
	ctx context.Context,
	r *connect.Request[participation_v1.RequestParticipationRequest],
) (*connect.Response[participation_v1.RequestParticipationResponse], error) {
	return x.requestParticipationRPC.RequestParticipation(ctx, r)
}

func (x Service) RestoreParticipation(
	ctx context.Context,
	r *connect.Request[participation_v1.RestoreParticipationRequest],
) (*connect.Response[emptypb.Empty], error) {
	return x.restoreParticipationRPC.RestoreParticipation(ctx, r)
}

func NewService(
	approveParticipationRPC rpc_service_participation_rpc_approve_participation.RPC,
	cancelParticipationRPC rpc_service_participation_rpc_cancel_participation.RPC,
	getParticipationRPC rpc_service_participation_rpc_get_participation.RPC,
	listParticipationChangesRPC rpc_service_participation_rpc_list_participation_changes.RPC,
	listParticipationsRPC rpc_service_participation_rpc_list_participations.RPC,
	rejectParticipationRPC rpc_service_participation_rpc_reject_participation.RPC,
	requestParticipationRPC rpc_service_participation_rpc_request_participation.RPC,
	restoreParticipationRPC rpc_service_participation_rpc_restore_participation.RPC,
) Service {
	return Service{
		approveParticipationRPC:     approveParticipationRPC,
		cancelParticipationRPC:      cancelParticipationRPC,
		getParticipationRPC:         getParticipationRPC,
		listParticipationChangesRPC: listParticipationChangesRPC,
		listParticipationsRPC:       listParticipationsRPC,
		rejectParticipationRPC:      rejectParticipationRPC,
		requestParticipationRPC:     requestParticipationRPC,
		restoreParticipationRPC:     restoreParticipationRPC,
	}
}
