package hahaton

import (
	"context"
	"flag"
	"fmt"

	"github.com/caarlos0/env/v10"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Roum1212/hahaton/internal/config"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

func RunSeed(ctx context.Context, arguments []string) error {
	flags := flag.NewFlagSet("seed", flag.ContinueOnError)
	coordinatorRaw := flags.String("coordinator-id", "", "existing coordinator user id")
	volunteerRaw := flags.String("volunteer-id", "", "existing volunteer user id")

	if err := flags.Parse(arguments); err != nil {
		return err
	}

	coordinatorID, err := model_user_id.NewUserID(*coordinatorRaw)
	if err != nil {
		return err
	}

	volunteerID, err := model_user_id.NewUserID(*volunteerRaw)
	if err != nil {
		return err
	}

	var cfg config.Config
	if err := env.ParseWithOptions(&cfg, config.ParseOptions()); err != nil {
		return fmt.Errorf("parse environment: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	pool, err := pgxpool.New(ctx, cfg.PostgreSQL.DSN)
	if err != nil {
		return err
	}
	defer pool.Close()

	key := cfg.DataEncryptionKey

	return Seed(ctx, pool, coordinatorID, volunteerID, key[:])
}
