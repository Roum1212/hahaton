package hahaton

import (
	"aidanwoods.dev/go-paseto"

	"github.com/Roum1212/hahaton/internal/config"
	paseto_rules "github.com/Roum1212/hahaton/internal/pkg/paseto"
	session_token "github.com/Roum1212/hahaton/internal/pkg/session/token"
	"github.com/Roum1212/hahaton/internal/pkg/token"
)

type tokenSet struct {
	sessionToken token.Service[session_token.Token]
}

func initTokens(cfg config.Config) tokenSet {
	return tokenSet{
		sessionToken: session_token.PASETO(
			cfg.SessionToken.V4SymmetricKey,
			paseto.MakeParser([]paseto.Rule{paseto_rules.ValidNow()}),
			cfg.SessionToken.TTL,
		),
	}
}
