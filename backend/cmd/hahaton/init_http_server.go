package hahaton

import (
	"net/http"
	"time"
)

const (
	headerTimeout = 5 * time.Second
	idleTimeout   = 60 * time.Second
)

func initHTTPServer(address string, handler http.Handler) *http.Server {
	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetHTTP2(true)
	protocols.SetUnencryptedHTTP2(true)

	return &http.Server{
		Protocols:         protocols,
		Addr:              address,
		Handler:           handler,
		ReadHeaderTimeout: headerTimeout,
		IdleTimeout:       idleTimeout,
	}
}
