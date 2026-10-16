package service_open_session

import (
	"context"
	"errors"
	"fmt"
	"time"

	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
	aggregate_session_command "github.com/Roum1212/hahaton/internal/domain/aggregate/session/command"
	model_session_id "github.com/Roum1212/hahaton/internal/domain/model/session-id"
	session_context "github.com/Roum1212/hahaton/internal/pkg/session/context"
	session_token "github.com/Roum1212/hahaton/internal/pkg/session/token"
	"github.com/Roum1212/hahaton/internal/pkg/token"
)

var ErrInvalidRequest = errors.New("invalid request")

type Service interface {
	OpenSession(ctx context.Context, r Request) (token.EncryptedToken, error)
}

type service struct {
	repository            aggregate_session.Repository
	sessionTokenEncryptor token.Encryptor[session_token.Token]
	ttl                   time.Duration
}

func (x service) OpenSession(
	ctx context.Context,
	r Request,
) (token.EncryptedToken, error) {
	if err := r.Validate(); err != nil {
		return token.EncryptedToken{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	if session, ok := session_context.GetSession(ctx); ok && session.IsActive(time.Now().UTC()) {
		return x.sessionTokenEncryptor.Encrypt(session_token.NewToken(
			session.GetSessionID(),
			session.GetUpdatedAt(),
			session.GetExpiresAt(),
		))
	}

	id, err := model_session_id.GenerateSessionID()
	if err != nil {
		return token.EncryptedToken{}, err
	}

	now := time.Now().UTC()
	command := aggregate_session_command.NewOpen(now, now.Add(x.ttl))
	session := aggregate_session.NewSession(id)

	if err = session.Open(command); err != nil {
		return token.EncryptedToken{}, err
	}

	encryptedToken, err := x.sessionTokenEncryptor.Encrypt(session_token.NewToken(
		session.GetSessionID(),
		session.GetUpdatedAt(),
		session.GetExpiresAt(),
	))
	if err != nil {
		return token.EncryptedToken{}, err
	}

	if err = x.repository.Create(ctx, session); err != nil {
		return token.EncryptedToken{}, fmt.Errorf("create session: %w", err)
	}

	return encryptedToken, nil
}

func NewService(
	repository aggregate_session.Repository,
	sessionTokenEncryptor token.Encryptor[session_token.Token],
	ttl time.Duration,
) Service {
	return service{
		repository:            repository,
		sessionTokenEncryptor: sessionTokenEncryptor,
		ttl:                   ttl,
	}
}
