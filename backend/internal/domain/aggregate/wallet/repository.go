package aggregate_wallet

import (
	"context"
	"errors"

	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	model_wallet_operation "github.com/Roum1212/hahaton/internal/domain/model/wallet-operation"
)

var ErrNotFound = errors.New("wallet not found")

type Repository interface {
	AppendOperation(ctx context.Context, value model_wallet_operation.Operation) error
	Create(ctx context.Context, wallet *Wallet) error
	ListOperations(ctx context.Context) ([]model_wallet_operation.Operation, error)
	Load(ctx context.Context, userID model_user_id.UserID) (*Wallet, error)
	Modify(
		ctx context.Context,
		userID model_user_id.UserID,
		change func(
			ctx context.Context,
			wallet *Wallet,
		) error,
	) error
	Save(ctx context.Context, wallet *Wallet) error
}
