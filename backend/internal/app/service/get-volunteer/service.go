package service_get_volunteer

import (
	"context"
	"errors"
	"fmt"

	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	aggregate_wallet "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrPermissionDenied = errors.New("permission denied")
)

type (
	Service interface {
		GetVolunteer(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		volunteers aggregate_volunteer.Repository
		wallets    aggregate_wallet.Repository
	}
)

func (x service) GetVolunteer(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	if value, ok := r.GetID().Get(); ok && value != r.GetUserID() {
		return Response{}, ErrPermissionDenied
	}

	volunteer, err := x.volunteers.Load(ctx, r.GetUserID())
	if err != nil {
		return Response{}, err
	}

	wallet, err := x.wallets.Load(ctx, r.GetUserID())
	if err != nil {
		return Response{}, err
	}

	return NewResponse(volunteer, wallet), nil
}

func NewService(volunteers aggregate_volunteer.Repository, wallets aggregate_wallet.Repository) Service {
	return service{
		volunteers: volunteers,
		wallets:    wallets,
	}
}
