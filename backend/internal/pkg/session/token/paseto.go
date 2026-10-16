package session_token

import (
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"

	model_session_id "github.com/Roum1212/hahaton/internal/domain/model/session-id"
	"github.com/Roum1212/hahaton/internal/pkg/token"
	paseto_token "github.com/Roum1212/hahaton/internal/pkg/token/paseto"
)

const (
	valueKey     = "v"
	sessionIDKey = "session_id"
)

func PASETO(key paseto.V4SymmetricKey, parser paseto.Parser, ttl time.Duration) token.Service[Token] {
	return paseto_token.NewService[Token](decrypt, encrypt, key, parser, ttl)
}

func decrypt(t paseto.Token) (Token, error) {
	var value map[string]string

	if err := t.Get(valueKey, &value); err != nil {
		return Token{}, fmt.Errorf("get session token value: %w", err)
	}

	sessionID, err := model_session_id.NewSessionID(value[sessionIDKey])
	if err != nil {
		return Token{}, fmt.Errorf("parse session token id: %w", err)
	}

	issuedAt, err := t.GetIssuedAt()
	if err != nil {
		return Token{}, fmt.Errorf("get session token issuance time: %w", err)
	}

	expiresAt, err := t.GetExpiration()
	if err != nil {
		return Token{}, fmt.Errorf("get session token expiration: %w", err)
	}

	return NewToken(sessionID, issuedAt, expiresAt), nil
}

func encrypt(t paseto.Token, value Token) error {
	if err := t.Set(valueKey, map[string]string{sessionIDKey: value.GetSessionID().String()}); err != nil {
		return fmt.Errorf("set session token value: %w", err)
	}

	t.SetExpiration(value.GetExpiresAt())
	t.SetIssuedAt(value.GetIssuedAt())
	t.SetNotBefore(value.GetIssuedAt())

	return nil
}
