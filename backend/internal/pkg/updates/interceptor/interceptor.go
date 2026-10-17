package updates_interceptor

import (
	"context"

	"connectrpc.com/connect"
)

type Interceptor struct {
	afterSuccess func(ctx context.Context, r connect.AnyRequest, response connect.AnyResponse)
}

func (x Interceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (x Interceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}

func (x Interceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, r connect.AnyRequest) (connect.AnyResponse, error) {
		response, err := next(ctx, r)
		if err == nil {
			x.afterSuccess(ctx, r, response)
		}

		return response, err
	}
}

func NewInterceptor(
	afterSuccess func(ctx context.Context, r connect.AnyRequest, response connect.AnyResponse),
) Interceptor {
	return Interceptor{afterSuccess: afterSuccess}
}
