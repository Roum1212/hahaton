package repository_event_postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	repository_locking "github.com/Roum1212/hahaton/internal/infra/repository/locking"
)

const (
	uniqueViolationCode = "23505"

	idColumn            = "id"
	coordinatorIDColumn = "coordinator_id"

	completionRequestIDColumn = "completion_request_id"
	descriptionColumn         = "description"
	imagePathColumn           = "image_path"
	latitudeColumn            = "latitude"
	locationColumn            = "location"
	longitudeColumn           = "longitude"
	requirementsColumn        = "requirements"
	resultFingerprintColumn   = "result_fingerprint"
	statusColumn              = "status"
	titleColumn               = "title"
	volunteerLimitColumn      = "volunteer_limit"

	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
	startsAtColumn  = "starts_at"
	endsAtColumn    = "ends_at"
)

type Events struct {
	database Database
	queries  squirrel.StatementBuilderType
}

func (x Events) Create(
	ctx context.Context,
	value *aggregate_event.Event,
	afterCreate func(
		ctx context.Context,
	) error,
) error {
	return x.database.InTx(ctx, func(ctx context.Context) error {
		if _, err := x.execute(ctx, x.queries.Insert("events").SetMap(aggregateValues(value))); err != nil {
			return err
		}

		return afterCreate(ctx)
	})
}

func (x Events) List(ctx context.Context) ([]*aggregate_event.Event, error) {
	sql, arguments, err := x.queries.Select(
		idColumn,
		coordinatorIDColumn,
		titleColumn,
		descriptionColumn,
		imagePathColumn,
		locationColumn,
		latitudeColumn,
		longitudeColumn,
		startsAtColumn,
		endsAtColumn,
		requirementsColumn,
		volunteerLimitColumn,
		statusColumn,
		"COALESCE("+completionRequestIDColumn+", '')",
		resultFingerprintColumn,
		createdAtColumn,
		updatedAtColumn,
	).From("events").ToSql()
	if err != nil {
		return nil, fmt.Errorf("build event list: %w", err)
	}

	rows, err := x.database.Query(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := make([]*aggregate_event.Event, 0)

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

func (x Events) Load(ctx context.Context, eventID model_event_id.EventID) (*aggregate_event.Event, error) {
	return x.load(ctx, eventID, false)
}

func (x Events) Modify(
	ctx context.Context,
	eventID model_event_id.EventID,
	change func(
		ctx context.Context,
		value *aggregate_event.Event,
	) error,
) error {
	return x.database.InTx(ctx, func(ctx context.Context) error {
		value, err := x.Load(ctx, eventID)
		if err != nil {
			return err
		}

		if err := repository_locking.Coordinator(ctx, x.database, value.GetCoordinatorID(), false); err != nil {
			return err
		}

		value, err = x.load(ctx, eventID, true)
		if err != nil {
			return err
		}

		if err := change(ctx, value); err != nil {
			return err
		}

		return x.Save(ctx, value)
	})
}

func (x Events) Save(ctx context.Context, value *aggregate_event.Event) error {
	values := aggregateValues(value)
	delete(values, idColumn)

	result, err := x.execute(
		ctx,
		x.queries.Update("events").SetMap(values).Where(squirrel.Eq{idColumn: value.GetEventID().String()}),
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return aggregate_event.ErrNotFound
	}

	return nil
}

func aggregateValues(value *aggregate_event.Event) map[string]any {
	imagePath := value.GetImagePath().String()

	source := value.GetRequirements().GetValues()

	requirements := make([]int32, len(source))
	for index, value := range source {
		requirements[index] = value.Int32()
	}

	var completionRequestID any
	if len(value.GetCompletionRequestID()) != 0 {
		completionRequestID = value.GetCompletionRequestID()
	}

	return map[string]any{
		idColumn:                  value.GetEventID().String(),
		coordinatorIDColumn:       value.GetCoordinatorID().String(),
		titleColumn:               value.GetTitle().String(),
		descriptionColumn:         value.GetDescription(),
		imagePathColumn:           imagePath,
		locationColumn:            value.GetLocation(),
		latitudeColumn:            value.GetPoint().GetLatitude(),
		longitudeColumn:           value.GetPoint().GetLongitude(),
		startsAtColumn:            value.GetSchedule().GetStartsAt(),
		endsAtColumn:              value.GetSchedule().GetEndsAt(),
		requirementsColumn:        requirements,
		volunteerLimitColumn:      value.GetVolunteerLimit().GetValue(),
		statusColumn:              value.GetStatus().String(),
		completionRequestIDColumn: completionRequestID,
		resultFingerprintColumn:   value.GetResultFingerprint(),
		createdAtColumn:           value.GetCreatedAt(),
		updatedAtColumn:           value.GetUpdatedAt(),
	}
}

func (x Events) execute(ctx context.Context, query squirrel.Sqlizer) (pgconn.CommandTag, error) {
	sql, arguments, err := query.ToSql()
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("build event query: %w", err)
	}

	result, err := x.database.Exec(ctx, sql, arguments...)
	if err != nil {
		var pgError *pgconn.PgError
		switch {
		case errors.As(err, &pgError) && pgError.Code == uniqueViolationCode:
			return pgconn.CommandTag{}, aggregate_event.ErrAlreadyExists
		default:
			return pgconn.CommandTag{}, fmt.Errorf("execute event query: %w", err)
		}
	}

	return result, nil
}

func (x Events) load(ctx context.Context, eventID model_event_id.EventID, lock bool) (*aggregate_event.Event, error) {
	query := x.queries.Select(
		idColumn,
		coordinatorIDColumn,
		titleColumn,
		descriptionColumn,
		imagePathColumn,
		locationColumn,
		latitudeColumn,
		longitudeColumn,
		startsAtColumn,
		endsAtColumn,
		requirementsColumn,
		volunteerLimitColumn,
		statusColumn,
		"COALESCE("+completionRequestIDColumn+", '')",
		resultFingerprintColumn,
		createdAtColumn,
		updatedAtColumn,
	).From("events").Where(squirrel.Eq{idColumn: eventID.String()})
	if lock {
		query = query.Suffix("FOR UPDATE")
	}

	sql, arguments, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var record row
	if err := record.scan(x.database.QueryRow(ctx, sql, arguments...)); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, aggregate_event.ErrNotFound
		default:
			return nil, err
		}
	}

	return record.toAggregate()
}

func NewEvents(database Database) Events {
	return Events{
		database: database,
		queries:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
