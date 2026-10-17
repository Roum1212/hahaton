//go:build e2e

package rpc

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	coordinator_v1 "github.com/Roum1212/hahaton/genproto/coordinator/v1"
	event_v1 "github.com/Roum1212/hahaton/genproto/event/v1"
	image_v1 "github.com/Roum1212/hahaton/genproto/image/v1"
	reward_v1 "github.com/Roum1212/hahaton/genproto/reward/v1"
	volunteer_v1 "github.com/Roum1212/hahaton/genproto/volunteer/v1"
)

func (x *RPCTestSuite) TestImagesOwnershipAndOptionalFields() {
	ctx := x.T().Context()
	c, outsider := x.coordinator(), x.login()

	var buffer bytes.Buffer
	x.Require().NoError(png.Encode(&buffer, image.NewRGBA(image.Rect(0, 0, 2, 2))))
	uploaded, err := x.clients.image.UploadImage(ctx, request(c, &image_v1.UploadImageRequest{Data: buffer.Bytes()}))
	x.Require().NoError(err)

	imagePath, imageURL := uploaded.Msg.GetImagePath(), uploaded.Msg.GetImageUrl()
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, x.baseURL+imageURL, nil)
	x.Require().NoError(err)
	response, err := x.httpClient.Do(r)
	x.Require().NoError(err)
	data, err := io.ReadAll(response.Body)
	x.Require().NoError(err)
	x.Require().NoError(response.Body.Close())
	x.Require().Equal(http.StatusOK, response.StatusCode)
	x.Require().Equal(buffer.Bytes(), data)
	_, err = x.clients.volunteer.UpdateVolunteer(ctx, request(outsider, &volunteer_v1.UpdateVolunteerRequest{
		AvatarPath: proto.String(imagePath),
	}))
	callError(x.T(), err, connect.CodePermissionDenied)
	_, err = x.clients.image.UploadImage(ctx, request(c, &image_v1.UploadImageRequest{Data: []byte("not an image")}))
	callError(x.T(), err, connect.CodeInvalidArgument)
	_, err = x.clients.image.UploadImage(ctx, request(c, &image_v1.UploadImageRequest{Data: make([]byte, 5<<20+1)}))
	callError(x.T(), err, connect.CodeInvalidArgument)
	_, err = x.clients.coordinator.UpdateCoordinator(ctx, request(c, &coordinator_v1.UpdateCoordinatorRequest{
		AvatarPath: proto.String(c.id + "/00000000-0000-7000-8000-000000000001.png"),
	}))
	callError(x.T(), err, connect.CodeNotFound)
	id := x.event(c, 1)

	rewardID := x.reward(c, 1, 1)
	for _, value := range []string{
		imagePath,
		"",
	} {
		_, err = x.clients.volunteer.UpdateVolunteer(ctx, request(c, &volunteer_v1.UpdateVolunteerRequest{
			AvatarPath: proto.String(value),
		}))
		x.Require().NoError(err)
		_, err = x.clients.coordinator.UpdateCoordinator(ctx, request(c, &coordinator_v1.UpdateCoordinatorRequest{
			AvatarPath: proto.String(value),
		}))
		x.Require().NoError(err)
		_, err = x.clients.event.UpdateEvent(ctx, request(c, &event_v1.UpdateEventRequest{
			Id:        id,
			ImagePath: proto.String(value),
		}))
		x.Require().NoError(err)
		_, err = x.clients.reward.UpdateReward(ctx, request(c, &reward_v1.UpdateRewardRequest{
			Id:        rewardID,
			ImagePath: proto.String(value),
		}))
		x.Require().NoError(err)
		_, err = x.clients.volunteer.UpdateVolunteer(ctx, request(c, &volunteer_v1.UpdateVolunteerRequest{}))
		x.Require().NoError(err)
		_, err = x.clients.coordinator.UpdateCoordinator(ctx, request(c, &coordinator_v1.UpdateCoordinatorRequest{}))
		x.Require().NoError(err)
		_, err = x.clients.event.UpdateEvent(ctx, request(c, &event_v1.UpdateEventRequest{Id: id}))
		x.Require().NoError(err)
		_, err = x.clients.reward.UpdateReward(ctx, request(c, &reward_v1.UpdateRewardRequest{Id: rewardID}))
		x.Require().NoError(err)
		profile, err := x.clients.volunteer.GetVolunteer(ctx, request(c, &volunteer_v1.GetVolunteerRequest{}))
		x.Require().NoError(err)
		coordinator, err := x.clients.coordinator.GetCoordinator(
			ctx,
			request(c, &coordinator_v1.GetCoordinatorRequest{}),
		)
		x.Require().NoError(err)
		event, err := x.clients.event.GetEvent(ctx, request(c, &event_v1.GetEventRequest{Id: id}))
		x.Require().NoError(err)
		reward, err := x.clients.reward.GetReward(ctx, request(c, &reward_v1.GetRewardRequest{Id: rewardID}))
		x.Require().NoError(err)

		expected := ""
		if value != "" {
			expected = imageURL
		}

		x.Require().Equal(expected, profile.Msg.GetVolunteer().GetAvatarUrl())
		x.Require().Equal(expected, coordinator.Msg.GetCoordinator().GetAvatarUrl())
		x.Require().Equal(expected, event.Msg.GetEvent().GetImageUrl())
		x.Require().Equal(expected, reward.Msg.GetReward().GetImageUrl())
	}
}
