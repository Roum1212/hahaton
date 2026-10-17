package repository_participation_postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	aggregate_participation "github.com/Roum1212/hahaton/internal/domain/aggregate/participation"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
)

const (
	uniqueViolationCode = "23505"

	idColumn              = "id"
	eventIDColumn         = "event_id"
	participationIDColumn = "participation_id"
	coordinatorIDColumn   = "coordinator_id"
	volunteerIDColumn     = "volunteer_id"

	accrualOperationIDColumn   = "accrual_operation_id"
	actorIDColumn              = "actor_id"
	actorKindColumn            = "actor_kind"
	awardedCoinsColumn         = "awarded_coins"
	awardedMinutesColumn       = "awarded_minutes"
	cancellationReasonColumn   = "cancellation_reason"
	hasChatAccessColumn        = "has_chat_access"
	previousStatusColumn       = "previous_status"
	rejectionReasonColumn      = "rejection_reason"
	sequenceColumn             = "sequence"
	statusColumn               = "status"
	visibleAfterSequenceColumn = "visible_after_sequence"

	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type Participations struct {
	database Database
	queries  squirrel.StatementBuilderType
}

func (x Participations) Create(ctx context.Context, value *aggregate_participation.Participation) error {
	_, err := x.execute(ctx, x.queries.Insert("participations").SetMap(aggregateValues(value)))

	return err
}

func (x Participations) List(ctx context.Context) ([]*aggregate_participation.Participation, error) {
	sql, arguments, err := x.queries.Select(
		idColumn,
		eventIDColumn,
		coordinatorIDColumn,
		volunteerIDColumn,
		statusColumn,
		cancellationReasonColumn,
		rejectionReasonColumn,
		awardedMinutesColumn,
		awardedCoinsColumn,
		accrualOperationIDColumn,
		hasChatAccessColumn,
		visibleAfterSequenceColumn,
		createdAtColumn,
		updatedAtColumn,
	).From("participations").ToSql()
	if err != nil {
		return nil, fmt.Errorf("build participation list: %w", err)
	}

	rows, err := x.database.Query(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := make([]*aggregate_participation.Participation, 0)

	for rows.Next() {
		var record row
		if err := record.scan(rows); err != nil {
			return nil, err
		}

		value, err := record.toAggregate()
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return values, nil
}

func (x Participations) Load(
	ctx context.Context,
	participationID model_participation_id.ParticipationID,
) (*aggregate_participation.Participation, error) {
	sql, arguments, err := x.queries.Select(
		idColumn,
		eventIDColumn,
		coordinatorIDColumn,
		volunteerIDColumn,
		statusColumn,
		cancellationReasonColumn,
		rejectionReasonColumn,
		awardedMinutesColumn,
		awardedCoinsColumn,
		accrualOperationIDColumn,
		hasChatAccessColumn,
		visibleAfterSequenceColumn,
		createdAtColumn,
		updatedAtColumn,
	).From("participations").Where(squirrel.Eq{idColumn: participationID.String()}).ToSql()
	if err != nil {
		return nil, err
	}

	var record row
	if err := record.scan(x.database.QueryRow(ctx, sql, arguments...)); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, aggregate_participation.ErrNotFound
		default:
			return nil, err
		}
	}

	return record.toAggregate()
}

func (x Participations) Save(ctx context.Context, value *aggregate_participation.Participation) error {
	values := aggregateValues(value)
	delete(values, idColumn)

	result, err := x.execute(
		ctx,
		x.queries.Update("participations").SetMap(values).Where(
			squirrel.Eq{idColumn: value.GetParticipationID().String()},
		),
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return aggregate_participation.ErrNotFound
	}

	return nil
}

func aggregateValues(value *aggregate_participation.Participation) map[string]any {
	operationID := ""
	if id := value.GetAccrualOperationID(); !id.IsZero() {
		operationID = id.String()
	}

	return map[string]any{
		idColumn:                   value.GetParticipationID().String(),
		eventIDColumn:              value.GetEventID().String(),
		coordinatorIDColumn:        value.GetCoordinatorID().String(),
		volunteerIDColumn:          value.GetVolunteerID().String(),
		statusColumn:               value.GetStatus().String(),
		cancellationReasonColumn:   value.GetCancellationReason(),
		rejectionReasonColumn:      value.GetRejectionReason(),
		awardedMinutesColumn:       value.GetAwardedMinutes(),
		awardedCoinsColumn:         value.GetAwardedCoins(),
		accrualOperationIDColumn:   operationID,
		hasChatAccessColumn:        value.GetHasChatAccess(),
		visibleAfterSequenceColumn: value.GetVisibleAfterSequence(),
		createdAtColumn:            value.GetCreatedAt(),
		updatedAtColumn:            value.GetUpdatedAt(),
	}
}

func (x Participations) execute(ctx context.Context, query squirrel.Sqlizer) (pgconn.CommandTag, error) {
	sql, arguments, err := query.ToSql()
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("build participation query: %w", err)
	}

	result, err := x.database.Exec(ctx, sql, arguments...)
	if err != nil {
		var pgError *pgconn.PgError
		switch {
		case errors.As(err, &pgError) && pgError.Code == uniqueViolationCode:
			return pgconn.CommandTag{}, aggregate_participation.ErrAlreadyExists
		default:
			return pgconn.CommandTag{}, fmt.Errorf("execute participation query: %w", err)
		}
	}

	return result, nil
}

func NewParticipations(database Database) Participations {
	return Participations{
		database: database,
		queries:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
