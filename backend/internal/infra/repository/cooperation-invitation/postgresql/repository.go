package repository_cooperation_invitation_postgresql

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	aggregate_cooperation_invitation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation"
	model_invitation_id "github.com/Roum1212/hahaton/internal/domain/model/invitation-id"
	repository_locking "github.com/Roum1212/hahaton/internal/infra/repository/locking"
)

const (
	uniqueViolationCode = "23505"

	idColumn            = "id"
	coordinatorIDColumn = "coordinator_id"

	tokenColumn = "token"

	createdAtColumn = "created_at"
	expiresAtColumn = "expires_at"
)

type Crypto interface {
	Decrypt(value []byte) (string, error)
	Encrypt(text string) ([]byte, error)
}

type CooperationInvitations struct {
	database Database
	crypto   Crypto
	queries  squirrel.StatementBuilderType
}

func (x CooperationInvitations) Create(
	ctx context.Context,
	value *aggregate_cooperation_invitation.CooperationInvitation,
) error {
	values, err := x.aggregateValues(value)
	if err != nil {
		return err
	}

	_, err = x.execute(ctx, x.queries.Insert("cooperation_invitations").SetMap(values))

	return err
}

func (x CooperationInvitations) Delete(ctx context.Context, invitationID model_invitation_id.InvitationID) error {
	_, err := x.execute(
		ctx,
		x.queries.Delete("cooperation_invitations").Where(squirrel.Eq{idColumn: invitationID.String()}),
	)

	return err
}

func (x CooperationInvitations) List(
	ctx context.Context,
) ([]*aggregate_cooperation_invitation.CooperationInvitation, error) {
	sql, arguments, err := x.queries.Select(
		idColumn,
		coordinatorIDColumn,
		tokenColumn,
		createdAtColumn,
		expiresAtColumn,
	).From(
		"cooperation_invitations",
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build cooperation-invitation list: %w", err)
	}

	rows, err := x.database.Query(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := make([]*aggregate_cooperation_invitation.CooperationInvitation, 0)

	for rows.Next() {
		var record row
		if err := record.scan(rows); err != nil {
			return nil, err
		}

		value, err := record.toAggregate(x.crypto)
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

func (x CooperationInvitations) Load(
	ctx context.Context,
	invitationID model_invitation_id.InvitationID,
) (*aggregate_cooperation_invitation.CooperationInvitation, error) {
	return x.load(ctx, invitationID, false)
}

func (x CooperationInvitations) Modify(
	ctx context.Context,
	invitationID model_invitation_id.InvitationID,
	change func(
		ctx context.Context,
		value *aggregate_cooperation_invitation.CooperationInvitation,
	) error,
) error {
	return x.database.InTx(ctx, func(ctx context.Context) error {
		value, err := x.Load(ctx, invitationID)
		if err != nil {
			return err
		}

		if err := repository_locking.Coordinator(ctx, x.database, value.GetCoordinatorID(), true); err != nil {
			return err
		}

		value, err = x.load(ctx, invitationID, true)
		if err != nil {
			return err
		}

		if err := change(ctx, value); err != nil {
			return err
		}

		return x.Save(ctx, value)
	})
}

func (x CooperationInvitations) Save(
	ctx context.Context,
	value *aggregate_cooperation_invitation.CooperationInvitation,
) error {
	values, err := x.aggregateValues(value)
	if err != nil {
		return err
	}

	delete(values, idColumn)

	result, err := x.execute(
		ctx,
		x.queries.Update("cooperation_invitations").SetMap(values).Where(squirrel.Eq{
			idColumn: value.GetInvitationID().String(),
		}),
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return aggregate_cooperation_invitation.ErrNotFound
	}

	return nil
}

func (x CooperationInvitations) aggregateValues(
	value *aggregate_cooperation_invitation.CooperationInvitation,
) (map[string]any, error) {
	encrypted, err := x.crypto.Encrypt(value.GetToken())
	if err != nil {
		return nil, err
	}

	return map[string]any{
		idColumn:            value.GetInvitationID().String(),
		coordinatorIDColumn: value.GetCoordinatorID().String(),
		tokenColumn:         base64.RawStdEncoding.EncodeToString(encrypted),
		createdAtColumn:     value.GetCreatedAt(),
		expiresAtColumn:     value.GetExpiresAt(),
	}, nil
}

func (x CooperationInvitations) execute(ctx context.Context, query squirrel.Sqlizer) (pgconn.CommandTag, error) {
	sql, arguments, err := query.ToSql()
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("build cooperation-invitation query: %w", err)
	}

	result, err := x.database.Exec(ctx, sql, arguments...)
	if err != nil {
		var pgError *pgconn.PgError
		switch {
		case errors.As(err, &pgError) && pgError.Code == uniqueViolationCode:
			return pgconn.CommandTag{}, aggregate_cooperation_invitation.ErrAlreadyExists
		default:
			return pgconn.CommandTag{}, fmt.Errorf("execute cooperation-invitation query: %w", err)
		}
	}

	return result, nil
}

func (x CooperationInvitations) load(
	ctx context.Context,
	invitationID model_invitation_id.InvitationID,
	lock bool,
) (*aggregate_cooperation_invitation.CooperationInvitation, error) {
	query := x.queries.Select(
		idColumn,
		coordinatorIDColumn,
		tokenColumn,
		createdAtColumn,
		expiresAtColumn,
	).From(
		"cooperation_invitations",
	).Where(
		squirrel.Eq{idColumn: invitationID.String()},
	)
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
			return nil, aggregate_cooperation_invitation.ErrNotFound
		default:
			return nil, err
		}
	}

	return record.toAggregate(x.crypto)
}

func NewCooperationInvitations(database Database, crypto Crypto) CooperationInvitations {
	return CooperationInvitations{
		database: database,
		crypto:   crypto,
		queries:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
