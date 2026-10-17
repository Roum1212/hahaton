package stream_controller

import (
	"context"
	"net/http"
)

func GetController(ctx context.Context) (*http.ResponseController, bool) {
	controller, ok := ctx.Value(controllerKey{}).(*http.ResponseController)

	return controller, ok
}

func WithController(ctx context.Context, controller *http.ResponseController) context.Context {
	return context.WithValue(ctx, controllerKey{}, controller)
}

type controllerKey struct{}
