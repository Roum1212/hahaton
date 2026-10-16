package rpc_service_session_serde

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	session_v1 "github.com/Roum1212/hahaton/genproto/session/v1"
	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
)

func AggregateToPBSession(session *aggregate_session.Session, hasCoordinatorProfile bool) *session_v1.Session {
	pb := &session_v1.Session{
		HasCoordinatorProfile: hasCoordinatorProfile,
		ExpiresAt:             timestamppb.New(session.GetExpiresAt()),
	}

	if !session.GetUserID().IsZero() {
		pb.UserId = proto.String(session.GetUserID().String())
	}

	return pb
}
