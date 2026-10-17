package rpc_serde

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func TimestampFromPB(value *timestamppb.Timestamp) (time.Time, error) {
	if err := value.CheckValid(); err != nil {
		return time.Time{}, err
	}

	return value.AsTime().UTC(), nil
}
