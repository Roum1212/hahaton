package repository_cooperation_postgresql

import (
	"fmt"
	"time"

	aggregate_cooperation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation"
	model_cooperation_status "github.com/Roum1212/hahaton/internal/domain/model/cooperation-status"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type row struct {
	coordinatorID string
	volunteerID   string
	status        string
	retryAfter    time.Time
	createdAt     time.Time
	updatedAt     time.Time
}

func (x *row) scan(source interface{ Scan(dest ...any) error }) error {
	return source.Scan(&x.coordinatorID, &x.volunteerID, &x.status, &x.retryAfter, &x.createdAt, &x.updatedAt)
}

func (x *row) toAggregate() (*aggregate_cooperation.Cooperation, error) {
	coordinatorID, err := model_user_id.NewUserID(x.coordinatorID)
	if err != nil {
		return nil, fmt.Errorf("restore coordinator id: %w", err)
	}

	volunteerID, err := model_user_id.NewUserID(x.volunteerID)
	if err != nil {
		return nil, fmt.Errorf("restore volunteer id: %w", err)
	}

	status, err := model_cooperation_status.NewCooperationStatus(x.status)
	if err != nil {
		return nil, fmt.Errorf("restore status: %w", err)
	}

	return aggregate_cooperation.NewCooperationFromFields(
		coordinatorID,
		volunteerID,
		status,
		x.retryAfter,
		x.createdAt,
		x.updatedAt,
	), nil
}
