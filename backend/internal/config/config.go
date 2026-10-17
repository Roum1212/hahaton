package config

import (
	"encoding/hex"
	"errors"
	"net/url"
	"reflect"

	"aidanwoods.dev/go-paseto"
	"github.com/caarlos0/env/v10"
)

type Config struct {
	AdminAPIKey       string        `env:"ADMIN_API_KEY,required,notEmpty"`
	AdminHTTPServer   AdminServer   `envPrefix:"ADMIN_HTTP_SERVER_"`
	CityCatalog       CityCatalog   `envPrefix:"CITY_CATALOG_"`
	DaData            DaData        `envPrefix:"DADATA_"`
	DataEncryptionKey EncryptionKey `env:"DATA_ENCRYPTION_KEY,required,notEmpty"`
	Flida             Flida         `envPrefix:"FLIDA_"`
	FrontendURL       string        `env:"FRONTEND_URL,notEmpty"                 envDefault:"https://localhost:3000"`
	HTTPServer        Server        `envPrefix:"HTTP_SERVER_"`
	PostgreSQL        PostgreSQL    `envPrefix:"POSTGRESQL_"`
	SessionToken      Token         `envPrefix:"SESSION_TOKEN_"`
	Storage           Storage       `envPrefix:"STORAGE_"`
}

func (x Config) Validate() error {
	if x.SessionToken.TTL <= 0 {
		return errors.New("session token ttl must be positive")
	}

	frontendURL, err := url.Parse(x.FrontendURL)
	if err != nil || frontendURL.Scheme != "https" || len(frontendURL.Host) == 0 {
		return errors.New("frontend url must be an absolute https url for secure session cookies")
	}

	if (len(x.HTTPServer.TLSCertificate) == 0) != (len(x.HTTPServer.TLSKey) == 0) {
		return errors.New("tls certificate and key must be provided together")
	}

	return nil
}

type AdminServer struct {
	Address string `env:"ADDRESS,notEmpty" envDefault:"127.0.0.1:8081"`
}

type CityCatalog struct {
	BoundariesPath string `env:"BOUNDARIES_PATH"`
	CSVPath        string `env:"CSV_PATH"         envDefault:"./data/city.csv"`
	Version        string `env:"VERSION,notEmpty" envDefault:"local"`
}

type DaData struct {
	APIKey string `env:"API_KEY"`
}

type Flida struct {
	ClientID string `env:"CLIENT_ID,required,notEmpty"`
	Issuer   string `env:"ISSUER,required,notEmpty"`
}

type PostgreSQL struct {
	DSN string `env:"DSN,required,notEmpty"`
}

type Server struct {
	Address        string `env:"ADDRESS,notEmpty" envDefault:":8080"`
	TLSCertificate string `env:"TLS_CERTIFICATE"`
	TLSKey         string `env:"TLS_KEY"`
}

type Storage struct {
	Directory string `env:"DIRECTORY,notEmpty" envDefault:"./storage"`
}

func ParseOptions() env.Options {
	return env.Options{
		FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeFor[EncryptionKey](): func(value string) (any, error) {
				decoded, err := hex.DecodeString(value)
				if err != nil || len(decoded) != 32 {
					return EncryptionKey{}, errors.New("data encryption key must contain 32 hex-encoded bytes")
				}

				return EncryptionKey(decoded), nil
			},
			reflect.TypeFor[paseto.V4SymmetricKey](): func(value string) (any, error) {
				return paseto.V4SymmetricKeyFromHex(value)
			},
		},
	}
}
