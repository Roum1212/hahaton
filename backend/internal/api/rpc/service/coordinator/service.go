package rpc_service_coordinator

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	coordinator_v1 "github.com/Roum1212/hahaton/genproto/coordinator/v1"
	coordinator_v1connect "github.com/Roum1212/hahaton/genproto/coordinator/v1/coordinator_v1connect"
	rpc_service_coordinator_rpc_get_coordinator "github.com/Roum1212/hahaton/internal/api/rpc/service/coordinator/rpc/get-coordinator"
	rpc_service_coordinator_rpc_get_coordinator_settings "github.com/Roum1212/hahaton/internal/api/rpc/service/coordinator/rpc/get-coordinator-settings"
	rpc_service_coordinator_rpc_list_coordinators "github.com/Roum1212/hahaton/internal/api/rpc/service/coordinator/rpc/list-coordinators"
	rpc_service_coordinator_rpc_update_coordinator "github.com/Roum1212/hahaton/internal/api/rpc/service/coordinator/rpc/update-coordinator"
	rpc_service_coordinator_rpc_update_coordinator_settings "github.com/Roum1212/hahaton/internal/api/rpc/service/coordinator/rpc/update-coordinator-settings"
)

var _ coordinator_v1connect.CoordinatorServiceHandler = Service{}

type Service struct {
	getCoordinatorRPC            rpc_service_coordinator_rpc_get_coordinator.RPC
	getCoordinatorSettingsRPC    rpc_service_coordinator_rpc_get_coordinator_settings.RPC
	listCoordinatorsRPC          rpc_service_coordinator_rpc_list_coordinators.RPC
	updateCoordinatorRPC         rpc_service_coordinator_rpc_update_coordinator.RPC
	updateCoordinatorSettingsRPC rpc_service_coordinator_rpc_update_coordinator_settings.RPC
}

func (x Service) GetCoordinator(
	ctx context.Context,
	r *connect.Request[coordinator_v1.GetCoordinatorRequest],
) (*connect.Response[coordinator_v1.GetCoordinatorResponse], error) {
	return x.getCoordinatorRPC.GetCoordinator(ctx, r)
}

func (x Service) GetCoordinatorSettings(
	ctx context.Context,
	r *connect.Request[coordinator_v1.GetCoordinatorSettingsRequest],
) (*connect.Response[coordinator_v1.GetCoordinatorSettingsResponse], error) {
	return x.getCoordinatorSettingsRPC.GetCoordinatorSettings(ctx, r)
}

func (x Service) ListCoordinators(
	ctx context.Context,
	r *connect.Request[coordinator_v1.ListCoordinatorsRequest],
) (*connect.Response[coordinator_v1.ListCoordinatorsResponse], error) {
	return x.listCoordinatorsRPC.ListCoordinators(ctx, r)
}

func (x Service) UpdateCoordinator(
	ctx context.Context,
	r *connect.Request[coordinator_v1.UpdateCoordinatorRequest],
) (*connect.Response[emptypb.Empty], error) {
	return x.updateCoordinatorRPC.UpdateCoordinator(ctx, r)
}

func (x Service) UpdateCoordinatorSettings(
	ctx context.Context,
	r *connect.Request[coordinator_v1.UpdateCoordinatorSettingsRequest],
) (*connect.Response[emptypb.Empty], error) {
	return x.updateCoordinatorSettingsRPC.UpdateCoordinatorSettings(ctx, r)
}

func NewService(
	getCoordinatorRPC rpc_service_coordinator_rpc_get_coordinator.RPC,
	getCoordinatorSettingsRPC rpc_service_coordinator_rpc_get_coordinator_settings.RPC,
	listCoordinatorsRPC rpc_service_coordinator_rpc_list_coordinators.RPC,
	updateCoordinatorRPC rpc_service_coordinator_rpc_update_coordinator.RPC,
	updateCoordinatorSettingsRPC rpc_service_coordinator_rpc_update_coordinator_settings.RPC,
) Service {
	return Service{
		getCoordinatorRPC:            getCoordinatorRPC,
		getCoordinatorSettingsRPC:    getCoordinatorSettingsRPC,
		listCoordinatorsRPC:          listCoordinatorsRPC,
		updateCoordinatorRPC:         updateCoordinatorRPC,
		updateCoordinatorSettingsRPC: updateCoordinatorSettingsRPC,
	}
}
