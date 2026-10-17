package repository_volunteer_postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	model_identity "github.com/Roum1212/hahaton/internal/domain/model/identity"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

func (x Volunteers) List(ctx context.Context) ([]*aggregate_volunteer.Volunteer, error) {
	sql, arguments, err := x.profileQuery().ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := x.database.Query(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := make([]*aggregate_volunteer.Volunteer, 0)

	for rows.Next() {
		value, err := scanProfile(rows)
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

func (x Volunteers) Load(ctx context.Context, userID model_user_id.UserID) (*aggregate_volunteer.Volunteer, error) {
	sql, arguments, err := x.profileQuery().Where(squirrel.Eq{idColumn: userID.String()}).ToSql()
	if err != nil {
		return nil, err
	}

	value, err := scanProfile(x.database.QueryRow(ctx, sql, arguments...))
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, aggregate_volunteer.ErrNotFound
		default:
			return nil, err
		}
	}

	return value, nil
}

func (x Volunteers) Modify(
	ctx context.Context,
	userID model_user_id.UserID,
	change func(
		ctx context.Context,
		volunteer *aggregate_volunteer.Volunteer,
	) error,
) error {
	return x.database.InTx(ctx, func(ctx context.Context) error {
		volunteer, err := x.Load(ctx, userID)
		if err != nil {
			return err
		}

		if err := change(ctx, volunteer); err != nil {
			return err
		}

		return x.Save(ctx, volunteer)
	})
}

func (x Volunteers) Save(ctx context.Context, volunteer *aggregate_volunteer.Volunteer) error {
	values := volunteerValues(volunteer)
	delete(values, idColumn)

	sql, arguments, err := x.queries.Update("volunteers").SetMap(values).Where(squirrel.Eq{
		idColumn: volunteer.GetUserID().String(),
	}).ToSql()
	if err != nil {
		return err
	}

	result, err := x.database.Exec(ctx, sql, arguments...)
	if err != nil {
		return fmt.Errorf("save volunteer: %w", err)
	}

	if result.RowsAffected() == 0 {
		return aggregate_volunteer.ErrNotFound
	}

	return nil
}

func (x Volunteers) profileQuery() squirrel.SelectBuilder {
	return x.queries.Select(
		idColumn,
		flidaIssuerColumn,
		flidaSubjectColumn,
		firstNameColumn,
		lastNameColumn,
		middleNameColumn,
		dateOfBirthColumn,
		emailAddressColumn,
		phoneNumberColumn,
		avatarPathColumn,
		createdAtColumn,
		updatedAtColumn,
	).From("volunteers")
}

func scanProfile(source interface{ Scan(dest ...any) error }) (*aggregate_volunteer.Volunteer, error) {
	var (
		row             volunteerRow
		issuer, subject string
	)

	if err := source.Scan(
		&row.id,
		&issuer,
		&subject,
		&row.firstName,
		&row.lastName,
		&row.middleName,
		&row.dateOfBirth,
		&row.emailAddress,
		&row.phoneNumber,
		&row.avatarPath,
		&row.createdAt,
		&row.updatedAt,
	); err != nil {
		return nil, err
	}

	identity, err := model_identity.NewIdentity(issuer, subject)
	if err != nil {
		return nil, fmt.Errorf("restore identity: %w", err)
	}

	return row.toVolunteer(identity)
}
