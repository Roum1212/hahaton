package repository_wallet_postgresql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	aggregate_wallet "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

const (
	idColumn            = "id"
	coordinatorIDColumn = "coordinator_id"
	volunteerIDColumn   = "volunteer_id"

	coinBalanceColumn           = "coin_balance"
	coinBalanceAfterColumn      = "coin_balance_after"
	coinDeltaColumn             = "coin_delta"
	confirmedMinutesColumn      = "confirmed_minutes"
	confirmedMinutesAfterColumn = "confirmed_minutes_after"
	descriptionColumn           = "description"
	lastSequenceColumn          = "last_sequence"
	minuteDeltaColumn           = "minute_delta"
	operationTypeColumn         = "operation_type"
	sequenceColumn              = "sequence"
	sourceIDColumn              = "source_id"

	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type Database interface {
	Exec(ctx context.Context, query string, arguments ...any) (pgconn.CommandTag, error)
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
	Query(ctx context.Context, query string, arguments ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, arguments ...any) pgx.Row
}

type Wallets struct {
	database Database
	queries  squirrel.StatementBuilderType
}

func (x Wallets) Create(ctx context.Context, wallet *aggregate_wallet.Wallet) error {
	query := x.queries.Insert("wallets").SetMap(map[string]any{
		volunteerIDColumn:      wallet.GetUserID().String(),
		coinBalanceColumn:      wallet.GetCoinBalance(),
		confirmedMinutesColumn: wallet.GetConfirmedMinutes(),
		lastSequenceColumn:     wallet.GetLastSequence(),
		createdAtColumn:        wallet.GetCreatedAt(),
		updatedAtColumn:        wallet.GetUpdatedAt(),
	})

	sql, arguments, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build wallet insert: %w", err)
	}

	if _, err := x.database.Exec(ctx, sql, arguments...); err != nil {
		return fmt.Errorf("create wallet: %w", err)
	}

	return nil
}

func (x Wallets) Load(ctx context.Context, userID model_user_id.UserID) (*aggregate_wallet.Wallet, error) {
	sql, arguments, err := x.queries.Select(
		coinBalanceColumn,
		confirmedMinutesColumn,
		lastSequenceColumn,
		createdAtColumn,
		updatedAtColumn,
	).From("wallets").Where(squirrel.Eq{volunteerIDColumn: userID.String()}).ToSql()
	if err != nil {
		return nil, err
	}

	var (
		balance, minutes     int64
		sequence             uint64
		createdAt, updatedAt time.Time
	)

	if err := x.database.QueryRow(ctx, sql, arguments...).Scan(
		&balance,
		&minutes,
		&sequence,
		&createdAt,
		&updatedAt,
	); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, aggregate_wallet.ErrNotFound
		default:
			return nil, err
		}
	}

	return aggregate_wallet.NewWalletFromFields(userID, balance, minutes, sequence, createdAt, updatedAt), nil
}

func (x Wallets) Modify(
	ctx context.Context,
	userID model_user_id.UserID,
	change func(
		ctx context.Context,
		wallet *aggregate_wallet.Wallet,
	) error,
) error {
	return x.database.InTx(ctx, func(ctx context.Context) error {
		wallet, err := x.Load(ctx, userID)
		if err != nil {
			return err
		}

		if err := change(ctx, wallet); err != nil {
			return err
		}

		return x.Save(ctx, wallet)
	})
}

func (x Wallets) Save(ctx context.Context, wallet *aggregate_wallet.Wallet) error {
	sql, arguments, err := x.queries.Update("wallets").SetMap(map[string]any{
		coinBalanceColumn:      wallet.GetCoinBalance(),
		confirmedMinutesColumn: wallet.GetConfirmedMinutes(),
		lastSequenceColumn:     wallet.GetLastSequence(),
		updatedAtColumn:        wallet.GetUpdatedAt(),
	}).Where(squirrel.Eq{volunteerIDColumn: wallet.GetUserID().String()}).ToSql()
	if err != nil {
		return err
	}

	result, err := x.database.Exec(ctx, sql, arguments...)
	if err != nil {
		return fmt.Errorf("save wallet: %w", err)
	}

	if result.RowsAffected() == 0 {
		return aggregate_wallet.ErrNotFound
	}

	return nil
}

func NewWallets(database Database) Wallets {
	return Wallets{
		database: database,
		queries:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
