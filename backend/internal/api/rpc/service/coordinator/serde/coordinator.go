package rpc_service_coordinator_serde

import (
	"strings"

	coordinator_v1 "github.com/Roum1212/hahaton/genproto/coordinator/v1"
	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
)

func AggregateToPBCoordinator(
	coordinator *aggregate_coordinator.Coordinator,
	volunteer *aggregate_volunteer.Volunteer,
) *coordinator_v1.Coordinator {
	name := volunteer.GetFullName()
	displayName := strings.TrimSpace(name.GetLastName() + " " + name.GetFirstName() + " " + name.GetMiddleName())

	return &coordinator_v1.Coordinator{
		Id:          coordinator.GetUserID().String(),
		AvatarUrl:   coordinator.GetAvatarPath().GetURL(),
		Description: coordinator.GetDescription(),
		DisplayName: displayName,
	}
}
