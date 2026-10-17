package hahaton

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Roum1212/hahaton/internal/config"
	client_postgresql "github.com/Roum1212/hahaton/internal/infra/client/postgresql"
)

const (
	shutdownTimeout = 5 * time.Second
)

func Run(ctx context.Context) error {
	var cfg config.Config
	if err := env.ParseWithOptions(&cfg, config.ParseOptions()); err != nil {
		return fmt.Errorf("parse environment: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	databaseClients, err := initDatabaseClients(ctx, cfg)
	if err != nil {
		return err
	}
	defer databaseClients.pool.Close()

	clientHandler, adminHandler, err := NewHandlers(ctx, cfg, databaseClients.pool)
	if err != nil {
		return err
	}

	clientServer := initHTTPServer(cfg.HTTPServer.Address, clientHandler)
	adminServer := initHTTPServer(cfg.AdminHTTPServer.Address, adminHandler)
	results := make(chan error, 2)

	go func() {
		if len(cfg.HTTPServer.TLSCertificate) != 0 {
			results <- clientServer.ListenAndServeTLS(cfg.HTTPServer.TLSCertificate, cfg.HTTPServer.TLSKey)

			return
		}

		results <- clientServer.ListenAndServe()
	}()
	go func() {
		results <- adminServer.ListenAndServe()
	}()

	var runError error
	select {
	case <-ctx.Done():
	case runError = <-results:
	}

	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
	defer cancel()

	clientError := clientServer.Shutdown(shutdownCtx)
	adminError := adminServer.Shutdown(shutdownCtx)

	if errors.Is(runError, http.ErrServerClosed) {
		runError = nil
	}

	return errors.Join(runError, clientError, adminError)
}

func NewHandlers(ctx context.Context, cfg config.Config, pool *pgxpool.Pool) (http.Handler, http.Handler, error) {
	externalClients, err := initExternalClients(ctx, cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("init external clients: %w", err)
	}

	databaseClients := databaseClientSet{
		pool:       pool,
		postgresql: client_postgresql.NewClient(pool),
	}

	repos, err := initRepositories(ctx, cfg, databaseClients, externalClients)
	if err != nil {
		return nil, nil, fmt.Errorf("init repositories: %w", err)
	}

	tokens := initTokens(cfg)

	supportServices := initSupportServices(cfg)
	services := initServices(cfg, externalClients, repos, tokens)
	rpcServices := initRPCServices(cfg, services, supportServices)
	clientHandler, adminHandler := initHTTPHandlers(cfg, repos, tokens, supportServices, rpcServices)

	return clientHandler, adminHandler, nil
}
