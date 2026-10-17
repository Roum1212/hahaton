//go:build e2e

package rpc

import (
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/emptypb"

	event_v1 "github.com/Roum1212/hahaton/genproto/event/v1"
	participation_v1 "github.com/Roum1212/hahaton/genproto/participation/v1"
	reward_order_v1 "github.com/Roum1212/hahaton/genproto/reward-order/v1"
	reward_v1 "github.com/Roum1212/hahaton/genproto/reward/v1"
)

func (x *RPCTestSuite) TestAuthenticationForEveryProtectedRPC() {
	ctx := x.T().Context()
	opened, err := x.clients.session.OpenSession(ctx, request(actor{}, &emptypb.Empty{}))
	x.Require().NoError(err)
	anonymous := actor{cookie: x.cookie(opened.Header())}
	closed := x.login()
	_, err = x.clients.session.CloseSession(ctx, request(closed, &emptypb.Empty{}))
	x.Require().NoError(err)

	invalid := actor{cookie: &http.Cookie{
		Name:  closed.cookie.Name,
		Value: "tampered",
	}}
	checked := 0

	protoregistry.GlobalFiles.RangeFiles(func(file protoreflect.FileDescriptor) bool {
		for i := 0; i < file.Services().Len(); i++ {
			service := file.Services().Get(i)
			if service.FullName() == "admin.v1.AdminService" || service.FullName() == "session.v1.SessionService" {
				continue
			}

			for j := 0; j < service.Methods().Len(); j++ {
				method := service.Methods().Get(j)
				if method.Name() == "ListCities" || method.Name() == "GetCity" {
					continue
				}

				procedure := "/" + string(service.FullName()) + "/" + string(method.Name())

				client := connect.NewClient[dynamicpb.Message, emptypb.Empty](
					x.httpClient,
					x.baseURL+procedure,
					connect.WithSchema(method),
				)
				for _, a := range []actor{
					{},
					anonymous,
					closed,
					invalid,
				} {
					r := request(a, dynamicpb.NewMessage(method.Input()))
					if method.IsStreamingServer() {
						stream, err := client.CallServerStream(ctx, r)
						if err == nil {
							x.Require().False(stream.Receive())
							err = stream.Err()
							x.Require().NoError(stream.Close())
						}

						callError(x.T(), err, connect.CodeUnauthenticated)

						continue
					}

					_, err := client.CallUnary(ctx, r)
					callError(x.T(), err, connect.CodeUnauthenticated)
				}

				checked++
			}
		}

		return true
	})
	x.Require().Equal(51, checked)
}

func (x *RPCTestSuite) TestListValidation() {
	c := x.coordinator()

	protoregistry.GlobalFiles.RangeFiles(func(file protoreflect.FileDescriptor) bool {
		for i := 0; i < file.Services().Len(); i++ {
			service := file.Services().Get(i)
			for j := 0; j < service.Methods().Len(); j++ {
				method := service.Methods().Get(j)
				if strings.HasPrefix(string(method.Name()), "List") {
					x.checkListInputs(c, method)
				}
			}
		}

		return true
	})
}

func (x *RPCTestSuite) TestNotFoundDetails() {
	ctx := x.T().Context()
	c := x.coordinator()

	const missingID = "00000000-0000-7000-8000-000000000001"

	cases := []struct {
		resource string
		call     func() error
	}{
		{
			resource: "event.v1.Event",
			call: func() error {
				_, err := x.clients.event.GetEvent(ctx, request(c, &event_v1.GetEventRequest{Id: missingID}))

				return err
			},
		},
		{
			resource: "reward.v1.Reward",
			call: func() error {
				_, err := x.clients.reward.GetReward(ctx, request(c, &reward_v1.GetRewardRequest{Id: missingID}))

				return err
			},
		},
		{
			resource: "participation.v1.Participation",
			call: func() error {
				_, err := x.clients.participation.GetParticipation(
					ctx,
					request(c, &participation_v1.GetParticipationRequest{Id: missingID}),
				)

				return err
			},
		},
		{
			resource: "reward_order.v1.RewardOrder",
			call: func() error {
				_, err := x.clients.rewardOrder.GetRewardOrder(
					ctx,
					request(c, &reward_order_v1.GetRewardOrderRequest{Id: missingID}),
				)

				return err
			},
		},
	}
	for _, scenario := range cases {
		err := scenario.call()
		callError(x.T(), err, connect.CodeNotFound)

		var rpcError *connect.Error
		x.Require().ErrorAs(err, &rpcError)
		x.Require().Len(rpcError.Details(), 1)
		detail, err := rpcError.Details()[0].Value()
		x.Require().NoError(err)

		info, ok := detail.(*errdetails.ResourceInfo)
		x.Require().True(ok)
		x.Require().Equal(scenario.resource, info.GetResourceType())
		x.Require().Contains(info.GetResourceName(), missingID)
		x.Require().NotEmpty(info.GetDescription())
	}
}

func (x *RPCTestSuite) checkListInputs(c actor, method protoreflect.MethodDescriptor) {
	const missingID = "00000000-0000-7000-8000-000000000001"

	fields := method.Input().Fields()
	if fields.ByName("limit") == nil || fields.ByName("offset") == nil {
		return
	}

	procedure := "/" + string(method.Parent().FullName()) + "/" + string(method.Name())

	client := connect.NewClient[dynamicpb.Message, emptypb.Empty](
		x.httpClient,
		x.baseURL+procedure,
		connect.WithSchema(method),
	)
	for name, value := range map[protoreflect.Name]protoreflect.Value{
		"limit":  protoreflect.ValueOfInt32(1001),
		"offset": protoreflect.ValueOfInt64(-1),
	} {
		message := dynamicpb.NewMessage(method.Input())
		if scope := fields.ByName("scope"); scope != nil {
			message.Set(scope, protoreflect.ValueOfEnum(2))
		}

		if eventID := fields.ByName("event_id"); eventID != nil {
			message.Set(eventID, protoreflect.ValueOfString(missingID))
		}

		if participationID := fields.ByName("participation_id"); participationID != nil {
			message.Set(participationID, protoreflect.ValueOfString(missingID))
		}

		message.Set(fields.ByName(name), value)
		_, err := client.CallUnary(x.T().Context(), request(c, message))
		callError(x.T(), err, connect.CodeInvalidArgument)
	}
}
