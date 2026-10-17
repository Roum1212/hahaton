package rpc_service_admin

import (
	"context"

	"connectrpc.com/connect"

	admin_v1 "github.com/Roum1212/hahaton/genproto/admin/v1"
	admin_v1connect "github.com/Roum1212/hahaton/genproto/admin/v1/admin_v1connect"
	rpc_service_admin_rpc_create_coordinator "github.com/Roum1212/hahaton/internal/api/rpc/service/admin/rpc/create-coordinator"
	rpc_service_admin_rpc_find_volunteers "github.com/Roum1212/hahaton/internal/api/rpc/service/admin/rpc/find-volunteers"
)

var _ admin_v1connect.AdminServiceHandler = Service{}

type Service struct {
	createCoordinatorRPC rpc_service_admin_rpc_create_coordinator.RPC
	findVolunteersRPC    rpc_service_admin_rpc_find_volunteers.RPC
}

func (x Service) CreateCoordinator(
	ctx context.Context,
	r *connect.Request[admin_v1.CreateCoordinatorRequest],
) (*connect.Response[admin_v1.CreateCoordinatorResponse], error) {
	return x.createCoordinatorRPC.CreateCoordinator(ctx, r)
}

func (x Service) FindVolunteers(
	ctx context.Context,
	r *connect.Request[admin_v1.FindVolunteersRequest],
) (*connect.Response[admin_v1.FindVolunteersResponse], error) {
	return x.findVolunteersRPC.FindVolunteers(ctx, r)
}

func NewService(
	createCoordinatorRPC rpc_service_admin_rpc_create_coordinator.RPC,
	findVolunteersRPC rpc_service_admin_rpc_find_volunteers.RPC,
) Service {
	return Service{
		createCoordinatorRPC: createCoordinatorRPC,
		findVolunteersRPC:    findVolunteersRPC,
	}
}
