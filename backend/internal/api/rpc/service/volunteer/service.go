package rpc_service_volunteer

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	volunteer_v1 "github.com/Roum1212/hahaton/genproto/volunteer/v1"
	volunteer_v1connect "github.com/Roum1212/hahaton/genproto/volunteer/v1/volunteer_v1connect"
	rpc_service_volunteer_rpc_get_coordinator_volunteer "github.com/Roum1212/hahaton/internal/api/rpc/service/volunteer/rpc/get-coordinator-volunteer"
	rpc_service_volunteer_rpc_get_volunteer "github.com/Roum1212/hahaton/internal/api/rpc/service/volunteer/rpc/get-volunteer"
	rpc_service_volunteer_rpc_list_volunteer_operations "github.com/Roum1212/hahaton/internal/api/rpc/service/volunteer/rpc/list-volunteer-operations"
	rpc_service_volunteer_rpc_update_volunteer "github.com/Roum1212/hahaton/internal/api/rpc/service/volunteer/rpc/update-volunteer"
)

var _ volunteer_v1connect.VolunteerServiceHandler = Service{}

type Service struct {
	getCoordinatorVolunteerRPC rpc_service_volunteer_rpc_get_coordinator_volunteer.RPC
	getVolunteerRPC            rpc_service_volunteer_rpc_get_volunteer.RPC
	listVolunteerOperationsRPC rpc_service_volunteer_rpc_list_volunteer_operations.RPC
	updateVolunteerRPC         rpc_service_volunteer_rpc_update_volunteer.RPC
}

func (x Service) GetCoordinatorVolunteer(
	ctx context.Context,
	r *connect.Request[volunteer_v1.GetCoordinatorVolunteerRequest],
) (*connect.Response[volunteer_v1.GetCoordinatorVolunteerResponse], error) {
	return x.getCoordinatorVolunteerRPC.GetCoordinatorVolunteer(ctx, r)
}

func (x Service) GetVolunteer(
	ctx context.Context,
	r *connect.Request[volunteer_v1.GetVolunteerRequest],
) (*connect.Response[volunteer_v1.GetVolunteerResponse], error) {
	return x.getVolunteerRPC.GetVolunteer(ctx, r)
}

func (x Service) ListVolunteerOperations(
	ctx context.Context,
	r *connect.Request[volunteer_v1.ListVolunteerOperationsRequest],
) (*connect.Response[volunteer_v1.ListVolunteerOperationsResponse], error) {
	return x.listVolunteerOperationsRPC.ListVolunteerOperations(ctx, r)
}

func (x Service) UpdateVolunteer(
	ctx context.Context,
	r *connect.Request[volunteer_v1.UpdateVolunteerRequest],
) (*connect.Response[emptypb.Empty], error) {
	return x.updateVolunteerRPC.UpdateVolunteer(ctx, r)
}

func NewService(
	getCoordinatorVolunteerRPC rpc_service_volunteer_rpc_get_coordinator_volunteer.RPC,
	getVolunteerRPC rpc_service_volunteer_rpc_get_volunteer.RPC,
	listVolunteerOperationsRPC rpc_service_volunteer_rpc_list_volunteer_operations.RPC,
	updateVolunteerRPC rpc_service_volunteer_rpc_update_volunteer.RPC,
) Service {
	return Service{
		getCoordinatorVolunteerRPC: getCoordinatorVolunteerRPC,
		getVolunteerRPC:            getVolunteerRPC,
		listVolunteerOperationsRPC: listVolunteerOperationsRPC,
		updateVolunteerRPC:         updateVolunteerRPC,
	}
}
