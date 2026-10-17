package repository_volunteer_postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	model_identity "github.com/Roum1212/hahaton/internal/domain/model/identity"
)

const (
	idColumn = "id"

	avatarPathColumn   = "avatar_path"
	dateOfBirthColumn  = "date_of_birth"
	emailAddressColumn = "email_address"
	firstNameColumn    = "first_name"
	flidaIssuerColumn  = "flida_issuer"
	flidaSubjectColumn = "flida_subject"
	lastNameColumn     = "last_name"
	middleNameColumn   = "middle_name"
	phoneNumberColumn  = "phone_number"

	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type Database interface {
	Exec(ctx context.Context, query string, arguments ...any) (pgconn.CommandTag, error)
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
	Query(ctx context.Context, query string, arguments ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, arguments ...any) pgx.Row
}

type Volunteers struct {
	database Database
	queries  squirrel.StatementBuilderType
}

func (x Volunteers) Create(
	ctx context.Context,
	volunteer *aggregate_volunteer.Volunteer,
	afterCreate func(ctx context.Context) error,
) error {
	return x.database.InTx(ctx, func(ctx context.Context) error {
		query := x.queries.Insert("volunteers").SetMap(volunteerValues(volunteer))

		sql, arguments, err := query.ToSql()
		if err != nil {
			return fmt.Errorf("build volunteer insert: %w", err)
		}

		if _, err := x.database.Exec(ctx, sql, arguments...); err != nil {
			return fmt.Errorf("create volunteer: %w", err)
		}

		return afterCreate(ctx)
	})
}

func (x Volunteers) FindByIdentity(
	ctx context.Context,
	identity model_identity.Identity,
) (*aggregate_volunteer.Volunteer, error) {
	query := x.queries.Select(
		idColumn, firstNameColumn, lastNameColumn, middleNameColumn, dateOfBirthColumn,
		emailAddressColumn, phoneNumberColumn, avatarPathColumn, createdAtColumn, updatedAtColumn,
	).From("volunteers").Where(squirrel.Eq{
		flidaIssuerColumn:  identity.GetIssuer(),
		flidaSubjectColumn: identity.GetSubject(),
	})

	sql, arguments, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build volunteer query: %w", err)
	}

	var row volunteerRow

	if err := x.database.QueryRow(ctx, sql, arguments...).Scan(
		&row.id, &row.firstName, &row.lastName, &row.middleName, &row.dateOfBirth,
		&row.emailAddress, &row.phoneNumber, &row.avatarPath, &row.createdAt, &row.updatedAt,
	); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, aggregate_volunteer.ErrNotFound
		default:
			return nil, fmt.Errorf("find volunteer: %w", err)
		}
	}

	return row.toVolunteer(identity)
}

func volunteerValues(volunteer *aggregate_volunteer.Volunteer) map[string]any {
	avatarPath := volunteer.GetAvatarPath().String()

	return map[string]any{
		idColumn:           volunteer.GetUserID().String(),
		flidaIssuerColumn:  volunteer.GetIdentity().GetIssuer(),
		flidaSubjectColumn: volunteer.GetIdentity().GetSubject(),
		firstNameColumn:    volunteer.GetFullName().GetFirstName(),
		lastNameColumn:     volunteer.GetFullName().GetLastName(),
		middleNameColumn:   volunteer.GetFullName().GetMiddleName(),
		dateOfBirthColumn:  volunteer.GetDateOfBirth().Time(),
		emailAddressColumn: volunteer.GetEmailAddress().String(),
		phoneNumberColumn:  volunteer.GetPhoneNumber().String(),
		avatarPathColumn:   avatarPath,
		createdAtColumn:    volunteer.GetCreatedAt(),
		updatedAtColumn:    volunteer.GetUpdatedAt(),
	}
}

func NewVolunteers(database Database) Volunteers {
	return Volunteers{
		database: database,
		queries:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
