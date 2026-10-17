package hahaton

import (
	"context"

	config "github.com/Roum1212/hahaton/internal/config"
	client_city_catalog "github.com/Roum1212/hahaton/internal/infra/client/city-catalog"
	repository_city_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/city/postgresql"
)

func initializeCities(ctx context.Context, cfg config.Config, catalog *repository_city_postgresql.Catalog) error {
	if len(catalog.GetVersion()) != 0 {
		return nil
	}

	source, err := client_city_catalog.NewClient(
		cfg.CityCatalog.CSVPath,
		cfg.CityCatalog.BoundariesPath,
		cfg.CityCatalog.Version,
	)
	if err != nil {
		return err
	}

	cities, err := source.ListCities(ctx)
	if err != nil {
		return err
	}

	return catalog.Save(ctx, cities, source.GetVersion())
}
