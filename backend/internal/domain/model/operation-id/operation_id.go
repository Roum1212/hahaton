package model_operation_id

import (
	"fmt"

	types_uuid_v7 "github.com/Roum1212/hahaton/internal/pkg/types/uuid-v7"
)

type OperationID struct {
	value types_uuid_v7.UUIDv7
}

func (x OperationID) IsZero() bool {
	return x.value.IsZero()
}

func (x OperationID) String() string {
	return x.value.String()
}

func (x OperationID) Validate() error {
	return x.value.Validate()
}

func GenerateOperationID() (OperationID, error) {
	id, err := types_uuid_v7.GenerateUUIDv7()
	if err != nil {
		return OperationID{}, fmt.Errorf("generate operation id: %w", err)
	}

	return NewOperationID(id.String())
}

func NewOperationID(value string) (OperationID, error) {
	id, err := types_uuid_v7.NewUUIDv7(value)
	if err != nil {
		return OperationID{}, fmt.Errorf("parse operation id: %w", err)
	}

	x := OperationID{value: id}
	if err := x.Validate(); err != nil {
		return OperationID{}, err
	}

	return x, nil
}
