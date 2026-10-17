package config

import (
	"time"

	"aidanwoods.dev/go-paseto"
)

type Token struct {
	TTL            time.Duration         `env:"TTL"                                envDefault:"24h"`
	V4SymmetricKey paseto.V4SymmetricKey `env:"V4_SYMMETRIC_KEY,required,notEmpty"`
}
