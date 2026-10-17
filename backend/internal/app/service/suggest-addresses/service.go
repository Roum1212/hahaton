package service_suggest_addresses

import (
	"context"
	"errors"
	"fmt"
	"strings"

	model_address_suggestion "github.com/Roum1212/hahaton/internal/domain/model/address-suggestion"
)

var (
	ErrUnavailable    = errors.New("address provider unavailable")
	ErrQuota          = errors.New("address provider quota exceeded")
	ErrInvalidRequest = errors.New("invalid request")
)

type (
	DaData interface {
		Suggest(
			ctx context.Context,
			query string,
			fromBound string,
			toBound string,
			count int32,
		) ([]model_address_suggestion.AddressSuggestion, error)
	}
	Service interface {
		SuggestAddresses(ctx context.Context, r Request) (Response, error)
	}
	service struct {
		dadata DaData
	}
)

func (x service) SuggestAddresses(ctx context.Context, r Request) (Response, error) {
	if err := r.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	count := r.GetCount()
	if count == 0 {
		count = 10
	}

	suggestions, err := x.dadata.Suggest(ctx, strings.TrimSpace(r.GetQuery()), r.GetFromBound(), r.GetToBound(), count)
	if err != nil {
		return Response{}, err
	}

	return NewResponse(suggestions), nil
}

func NewService(dadata DaData) Service {
	return service{
		dadata: dadata,
	}
}
