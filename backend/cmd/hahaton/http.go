package hahaton

import (
	"crypto/subtle"
	"net/http"

	connectcors "connectrpc.com/cors"
	"github.com/rs/cors"

	stream_controller "github.com/Roum1212/hahaton/internal/pkg/stream/controller"
)

func withAdminAuthentication(next http.Handler, key string) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		expected := "Bearer " + key
		if len(key) == 0 || subtle.ConstantTimeCompare([]byte(r.Header.Get("Authorization")), []byte(expected)) != 1 {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func withCORS(next http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   append(connectcors.AllowedHeaders(), "Authorization"),
		AllowedMethods:   connectcors.AllowedMethods(),
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		ExposedHeaders: connectcors.ExposedHeaders(),
	}).Handler(next)
}

func withStreamingController(next http.Handler) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		ctx := stream_controller.WithController(r.Context(), http.NewResponseController(w))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
