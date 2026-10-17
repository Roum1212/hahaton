package rpc_error

import (
	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/proto"
)

func NewAbortedError(details *errdetails.ErrorInfo) error {
	return newError(connect.CodeAborted, details)
}

func NewAlreadyExistsError(details *errdetails.ResourceInfo) error {
	return newError(connect.CodeAlreadyExists, details)
}

func NewCanceledError() error {
	return newError(connect.CodeCanceled, nil)
}

func NewDataLossError() error {
	return newError(connect.CodeDataLoss, nil)
}

func NewDeadlineExceededError() error {
	return newError(connect.CodeDeadlineExceeded, nil)
}

func NewFailedPreconditionError(details *errdetails.PreconditionFailure) error {
	return newError(connect.CodeFailedPrecondition, details)
}

func NewInternalError() error {
	return newError(connect.CodeInternal, nil)
}

func NewInvalidArgumentError(details *errdetails.BadRequest) error {
	return newError(connect.CodeInvalidArgument, details)
}

func NewNotFoundError(details *errdetails.ResourceInfo) error {
	return newError(connect.CodeNotFound, details)
}

func NewOutOfRangeError(details *errdetails.BadRequest) error {
	return newError(connect.CodeOutOfRange, details)
}

func NewPermissionDeniedError(details *errdetails.ErrorInfo) error {
	return newError(connect.CodePermissionDenied, details)
}

func NewResourceExhaustedError(details *errdetails.QuotaFailure) error {
	return newError(connect.CodeResourceExhausted, details)
}

func NewUnauthenticatedError(details *errdetails.ErrorInfo) error {
	return newError(connect.CodeUnauthenticated, details)
}

func NewUnavailableError() error {
	return newError(connect.CodeUnavailable, nil)
}

func NewUnimplementedError() error {
	return newError(connect.CodeUnimplemented, nil)
}

func NewUnknownError() error {
	return newError(connect.CodeUnknown, nil)
}

func newError(code connect.Code, m proto.Message) error {
	e := connect.NewError(code, nil)

	d, err := connect.NewErrorDetail(m)
	if err != nil {
		return e
	}

	e.AddDetail(d)

	return e
}
