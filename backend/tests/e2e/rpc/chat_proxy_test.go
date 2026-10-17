//go:build e2e

package rpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	event_chat_v1 "github.com/Roum1212/hahaton/genproto/event-chat/v1"
	event_chat_v1connect "github.com/Roum1212/hahaton/genproto/event-chat/v1/event_chat_v1connect"
)

const streamDeliveryTimeout = 3 * time.Second

func (x *RPCTestSuite) TestChatStreamingThroughBufferingProxy() {
	ctx := x.T().Context()
	coordinator, volunteer := x.coordinator(), x.login()
	x.cooperate(coordinator, volunteer)
	eventID := x.event(coordinator, 1)
	x.participate(coordinator, volunteer, eventID)

	_, err := x.clients.eventChat.SendEventChatMessage(ctx, request(coordinator,
		&event_chat_v1.SendEventChatMessageRequest{EventId: eventID, Text: "Первое", RequestId: "proxy-first"}))
	x.Require().NoError(err)

	proxyURL := x.bufferingProxy()
	client := event_chat_v1connect.NewEventChatServiceClient(&http.Client{Timeout: requestTimeout}, proxyURL)

	streamCtx, cancel := context.WithTimeout(ctx, streamDeliveryTimeout)
	defer cancel()

	stream, err := client.SubscribeEventChatMessages(streamCtx,
		request(coordinator, &event_chat_v1.SubscribeEventChatMessagesRequest{EventId: eventID}))
	x.Require().NoError(err)
	x.T().Cleanup(func() { x.Require().NoError(stream.Close()) })
	x.Require().True(stream.Receive(), "first message is buffered: %v", stream.Err())
	x.Require().Equal("Первое", stream.Msg().GetMessage().GetText())

	_, err = x.clients.eventChat.SendEventChatMessage(ctx, request(volunteer,
		&event_chat_v1.SendEventChatMessageRequest{EventId: eventID, Text: "Новое", RequestId: "proxy-live"}))
	x.Require().NoError(err)
	x.Require().True(stream.Receive(), "new message is buffered: %v", stream.Err())
	x.Require().Equal("Новое", stream.Msg().GetMessage().GetText())
}

func (x *RPCTestSuite) bufferingProxy() string {
	ctx := x.T().Context()
	upstream, err := url.Parse(x.baseURL)
	x.Require().NoError(err)
	port, err := strconv.Atoi(upstream.Port())
	x.Require().NoError(err)
	configPath := filepath.Join(x.T().TempDir(), "default.conf")
	config := fmt.Sprintf(`server {
    listen 80;
    location / {
        proxy_pass https://host.testcontainers.internal:%d;
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 4 4k;
    }
}`, port)
	x.Require().NoError(os.WriteFile(configPath, []byte(config), 0o600))
	container, err := testcontainers.Run(ctx, "nginx:1.28-alpine",
		testcontainers.WithHostPortAccess(port),
		testcontainers.WithFiles(testcontainers.ContainerFile{
			HostFilePath:      configPath,
			ContainerFilePath: "/etc/nginx/conf.d/default.conf",
			FileMode:          0o644,
		}),
		testcontainers.WithExposedPorts("80/tcp"),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("80/tcp")),
	)
	x.Require().NoError(err)
	testcontainers.CleanupContainer(x.T(), container)
	host, err := container.Host(ctx)
	x.Require().NoError(err)
	proxyPort, err := container.MappedPort(ctx, "80/tcp")
	x.Require().NoError(err)

	return "http://" + net.JoinHostPort(host, proxyPort.Port())
}
