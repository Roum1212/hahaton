package repository_cooperation_postgresql

import (
	"context"

	model_cooperation_change "github.com/Roum1212/hahaton/internal/domain/model/cooperation-change"
)

func (x Cooperations) AppendChange(ctx context.Context, change model_cooperation_change.CooperationChange) error {
	values := map[string]any{
		coordinatorIDColumn: change.GetCoordinatorID().String(),
		volunteerIDColumn:   change.GetVolunteerID().String(),
		statusColumn:        change.GetStatus().String(),
		actorIDColumn:       change.GetActorID().String(),
		retryAfterColumn:    nil,
		createdAtColumn:     change.GetCreatedAt(),
	}
	if !change.GetRetryAfter().IsZero() {
		values[retryAfterColumn] = change.GetRetryAfter()
	}

	query, arguments, err := x.queries.Insert("cooperation_changes").SetMap(values).ToSql()
	if err != nil {
		return err
	}

	_, err = x.database.Exec(ctx, query, arguments...)

	return err
}
