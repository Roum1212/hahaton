package service_list_volunteer_operations

import (
	"fmt"

	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	model_wallet_operation "github.com/Roum1212/hahaton/internal/domain/model/wallet-operation"
	model_wallet_operation_type "github.com/Roum1212/hahaton/internal/domain/model/wallet-operation-type"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

type Request struct {
	userID model_user_id.UserID

	operationType utils_value.Optional[model_wallet_operation_type.OperationType]

	page model_page.Page
}

func (x Request) GetOperationType() utils_value.Optional[model_wallet_operation_type.OperationType] {
	return x.operationType
}

func (x Request) GetPage() model_page.Page {
	return x.page
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if err := x.page.Validate(); err != nil {
		return fmt.Errorf("page: %w", err)
	}

	if value, ok := x.operationType.Get(); ok {
		return value.Validate()
	}

	return nil
}

type RequestOption func(x *Request)

func WithOperationType(value model_wallet_operation_type.OperationType) RequestOption {
	return func(x *Request) {
		x.operationType = utils_value.NewOptional(value)
	}
}

func NewRequest(userID model_user_id.UserID, page model_page.Page, opts ...RequestOption) Request {
	x := Request{
		userID: userID,
		page:   page,
	}
	for _, option := range opts {
		option(&x)
	}

	return x
}

type Response struct {
	operations []model_wallet_operation.Operation
	totalCount int64
}

func (x Response) GetOperations() []model_wallet_operation.Operation {
	return x.operations
}

func (x Response) GetTotalCount() int64 {
	return x.totalCount
}

func NewResponse(operations []model_wallet_operation.Operation, totalCount int64) Response {
	return Response{
		operations: operations,
		totalCount: totalCount,
	}
}
