package service_list_volunteer_operations

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"

	aggregate_wallet "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet"
	model_wallet_operation "github.com/Roum1212/hahaton/internal/domain/model/wallet-operation"
	"github.com/Roum1212/hahaton/internal/pkg/page"
)

var ErrInvalidRequest = errors.New("invalid request")

type (
	Service interface {
		ListVolunteerOperations(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		wallets aggregate_wallet.Repository
	}
)

func (x service) ListVolunteerOperations(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	operations, err := x.wallets.ListOperations(ctx)
	if err != nil {
		return Response{}, err
	}

	filtered := make([]model_wallet_operation.Operation, 0)

	for _, value := range operations {
		if value.GetVolunteerID() != r.GetUserID() {
			continue
		}

		if filter, ok := r.GetOperationType().Get(); ok && value.GetOperationType() != filter {
			continue
		}

		filtered = append(filtered, value)
	}

	slices.SortFunc(filtered, func(
		a model_wallet_operation.Operation,
		b model_wallet_operation.Operation,
	) int {
		return cmp.Compare(b.GetSequence(), a.GetSequence())
	})
	values, total := page.Slice(filtered, r.GetPage())

	return NewResponse(values, total), nil
}

func NewService(wallets aggregate_wallet.Repository) Service {
	return service{
		wallets: wallets,
	}
}
