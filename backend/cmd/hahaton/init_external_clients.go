package hahaton

import (
	"context"
	"net/http"
	"time"

	"github.com/Roum1212/hahaton/internal/config"
	client_dadata "github.com/Roum1212/hahaton/internal/infra/client/dadata"
	client_encryption "github.com/Roum1212/hahaton/internal/infra/client/encryption"
	client_flida "github.com/Roum1212/hahaton/internal/infra/client/flida"
	client_image_storage "github.com/Roum1212/hahaton/internal/infra/client/image-storage"
)

const providerTimeout = 10 * time.Second

type externalClientSet struct {
	dadata     client_dadata.Client
	encryption client_encryption.Client
	flida      client_flida.Client
	images     client_image_storage.Client
}

func initExternalClients(ctx context.Context, cfg config.Config) (externalClientSet, error) {
	httpClient := &http.Client{Timeout: providerTimeout}

	dadata, err := client_dadata.NewClient(cfg.DaData.APIKey, cfg.DaData.Endpoint, httpClient)
	if err != nil {
		return externalClientSet{}, err
	}

	flida, err := client_flida.NewClient(ctx, cfg.Flida.Issuer, cfg.Flida.ClientID, httpClient)
	if err != nil {
		return externalClientSet{}, err
	}

	key := cfg.DataEncryptionKey

	encryption, err := client_encryption.NewClient(key[:])
	if err != nil {
		return externalClientSet{}, err
	}

	return externalClientSet{
		dadata:     dadata,
		encryption: encryption,
		flida:      flida,
		images:     client_image_storage.NewClient(cfg.Storage.Directory),
	}, nil
}
