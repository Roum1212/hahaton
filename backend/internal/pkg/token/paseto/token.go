package paseto_token

import (
	"fmt"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"

	"github.com/Roum1212/hahaton/internal/pkg/token"
)

var prefix = paseto.V4Local.Header()

type Decrypt[T any] func(t paseto.Token) (T, error)

type Encrypt[T any] func(
	t paseto.Token,
	value T,
) error

type service[T any] struct {
	decrypt Decrypt[T]
	encrypt Encrypt[T]
	key     paseto.V4SymmetricKey
	parser  paseto.Parser
	ttl     time.Duration
}

func (x service[T]) Decrypt(encryptedToken token.EncryptedToken) (T, error) {
	t, err := x.parser.ParseV4Local(x.key, prefix+encryptedToken.String(), nil)
	if err != nil {
		var zero T

		return zero, fmt.Errorf("parse paseto: %w", err)
	}

	return x.decrypt(*t)
}

func (x service[T]) Encrypt(value T) (token.EncryptedToken, error) {
	t := paseto.NewToken()
	now := time.Now().UTC()

	t.SetExpiration(now.Add(x.ttl))
	t.SetIssuedAt(now)
	t.SetNotBefore(now)

	if err := x.encrypt(t, value); err != nil {
		return token.EncryptedToken{}, fmt.Errorf("encrypt paseto: %w", err)
	}

	encrypted := t.V4Encrypt(x.key, nil)

	return token.NewEncryptedToken(strings.TrimPrefix(encrypted, prefix))
}

func NewService[T any](
	decrypt Decrypt[T],
	encrypt Encrypt[T],
	key paseto.V4SymmetricKey,
	parser paseto.Parser,
	ttl time.Duration,
) token.Service[T] {
	return service[T]{
		decrypt: decrypt,
		encrypt: encrypt,
		key:     key,
		parser:  parser,
		ttl:     ttl,
	}
}
