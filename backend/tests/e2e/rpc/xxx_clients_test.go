//go:build e2e

package rpc

import (
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	admin_v1connect "github.com/Roum1212/hahaton/genproto/admin/v1/admin_v1connect"
	cooperation_v1connect "github.com/Roum1212/hahaton/genproto/cooperation/v1/cooperation_v1connect"
	coordinator_v1connect "github.com/Roum1212/hahaton/genproto/coordinator/v1/coordinator_v1connect"
	event_chat_v1connect "github.com/Roum1212/hahaton/genproto/event-chat/v1/event_chat_v1connect"
	event_v1connect "github.com/Roum1212/hahaton/genproto/event/v1/event_v1connect"
	geo_v1connect "github.com/Roum1212/hahaton/genproto/geo/v1/geo_v1connect"
	image_v1connect "github.com/Roum1212/hahaton/genproto/image/v1/image_v1connect"
	participation_v1connect "github.com/Roum1212/hahaton/genproto/participation/v1/participation_v1connect"
	reward_order_v1connect "github.com/Roum1212/hahaton/genproto/reward-order/v1/reward_order_v1connect"
	reward_v1connect "github.com/Roum1212/hahaton/genproto/reward/v1/reward_v1connect"
	session_v1connect "github.com/Roum1212/hahaton/genproto/session/v1/session_v1connect"
	updates_v1connect "github.com/Roum1212/hahaton/genproto/updates/v1/updates_v1connect"
	volunteer_v1connect "github.com/Roum1212/hahaton/genproto/volunteer/v1/volunteer_v1connect"
)

type rpcClients struct {
	updates       updates_v1connect.UpdatesServiceClient
	admin         admin_v1connect.AdminServiceClient
	cooperation   cooperation_v1connect.CooperationServiceClient
	coordinator   coordinator_v1connect.CoordinatorServiceClient
	event         event_v1connect.EventServiceClient
	eventChat     event_chat_v1connect.EventChatServiceClient
	geo           geo_v1connect.GeoServiceClient
	image         image_v1connect.ImageServiceClient
	participation participation_v1connect.ParticipationServiceClient
	reward        reward_v1connect.RewardServiceClient
	rewardOrder   reward_order_v1connect.RewardOrderServiceClient
	session       session_v1connect.SessionServiceClient
	volunteer     volunteer_v1connect.VolunteerServiceClient
}

func newRPCClients(client *http.Client, baseURL string) rpcClients {
	return rpcClients{
		updates:       updates_v1connect.NewUpdatesServiceClient(client, baseURL),
		admin:         admin_v1connect.NewAdminServiceClient(client, baseURL),
		cooperation:   cooperation_v1connect.NewCooperationServiceClient(client, baseURL),
		coordinator:   coordinator_v1connect.NewCoordinatorServiceClient(client, baseURL),
		event:         event_v1connect.NewEventServiceClient(client, baseURL),
		eventChat:     event_chat_v1connect.NewEventChatServiceClient(client, baseURL),
		geo:           geo_v1connect.NewGeoServiceClient(client, baseURL),
		image:         image_v1connect.NewImageServiceClient(client, baseURL),
		participation: participation_v1connect.NewParticipationServiceClient(client, baseURL),
		reward:        reward_v1connect.NewRewardServiceClient(client, baseURL),
		rewardOrder:   reward_order_v1connect.NewRewardOrderServiceClient(client, baseURL),
		session:       session_v1connect.NewSessionServiceClient(client, baseURL),
		volunteer:     volunteer_v1connect.NewVolunteerServiceClient(client, baseURL),
	}
}

type observedTransport struct {
	next     http.RoundTripper
	coverage *rpcCoverage
}

func (x *observedTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	response, err := x.next.RoundTrip(r)
	if err == nil && response.StatusCode == http.StatusOK {
		x.coverage.mu.Lock()
		x.coverage.success[r.URL.Path] = true
		x.coverage.mu.Unlock()
	}

	return response, err
}

type rpcCoverage struct {
	mu      sync.Mutex
	success map[string]bool
}

func (x *rpcCoverage) assertComplete(t *testing.T) {
	t.Helper()
	x.mu.Lock()
	defer x.mu.Unlock()

	count := 0

	protoregistry.GlobalFiles.RangeFiles(func(file protoreflect.FileDescriptor) bool {
		if file.Services().Len() == 0 {
			return true
		}

		services := file.Services()
		for i := 0; i < services.Len(); i++ {
			service := services.Get(i)
			for j := 0; j < service.Methods().Len(); j++ {
				method := service.Methods().Get(j)
				procedure := "/" + string(service.FullName()) + "/" + string(method.Name())
				assert.True(t, x.success[procedure], "RPC has no successful end-to-end call: %s", procedure)

				count++
			}
		}

		return true
	})
	require.Equal(t, 60, count)
	t.Logf("verified successful calls of all %d RPC procedures", count)
}

func newRPCCoverage() *rpcCoverage {
	return &rpcCoverage{success: make(map[string]bool)}
}
