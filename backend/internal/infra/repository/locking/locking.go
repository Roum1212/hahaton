package repository_locking

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_reward_id "github.com/Roum1212/hahaton/internal/domain/model/reward-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

const (
	idColumn            = "id"
	coordinatorIDColumn = "coordinator_id"
)

type Database interface {
	QueryRow(ctx context.Context, query string, arguments ...any) pgx.Row
}

func Coordinator(ctx context.Context, database Database, id model_user_id.UserID, exclusive bool) error {
	mode := "FOR SHARE"
	if exclusive {
		mode = "FOR UPDATE"
	}

	query := squirrel.Select(idColumn).From("coordinators").Where(squirrel.Eq{idColumn: id.String()}).Suffix(mode)

	return acquire(ctx, database, query, aggregate_coordinator.ErrNotFound)
}

func Event(ctx context.Context, database Database, id model_event_id.EventID) error {
	query := squirrel.Select(coordinatorIDColumn).From("events").Where(squirrel.Eq{idColumn: id.String()})

	sql, arguments, err := query.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	var owner string
	if err := database.QueryRow(ctx, sql, arguments...).Scan(&owner); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return aggregate_event.ErrNotFound
		default:
			return fmt.Errorf("load event owner: %w", err)
		}
	}

	coordinatorID, err := model_user_id.NewUserID(owner)
	if err != nil {
		return err
	}

	if err := Coordinator(ctx, database, coordinatorID, false); err != nil {
		return err
	}

	query = squirrel.Select(idColumn).From("events").Where(squirrel.Eq{idColumn: id.String()}).Suffix("FOR UPDATE")

	return acquire(ctx, database, query, aggregate_event.ErrNotFound)
}

func Reward(ctx context.Context, database Database, id model_reward_id.RewardID) error {
	query := squirrel.Select(idColumn).From("rewards").Where(squirrel.Eq{idColumn: id.String()}).Suffix("FOR UPDATE")

	return acquire(ctx, database, query, aggregate_reward.ErrNotFound)
}

func acquire(ctx context.Context, database Database, query squirrel.SelectBuilder, notFound error) error {
	sql, arguments, err := query.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	var id string
	if err := database.QueryRow(ctx, sql, arguments...).Scan(&id); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return notFound
		default:
			return fmt.Errorf("lock row: %w", err)
		}
	}

	return nil
}
