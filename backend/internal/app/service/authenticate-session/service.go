package service_authenticate_session

import (
	"context"
	"errors"
	"fmt"
	"time"

	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
	aggregate_session_command "github.com/Roum1212/hahaton/internal/domain/aggregate/session/command"
	model_session_id "github.com/Roum1212/hahaton/internal/domain/model/session-id"
	session_token "github.com/Roum1212/hahaton/internal/pkg/session/token"
	"github.com/Roum1212/hahaton/internal/pkg/token"
)

var ErrInvalidRequest = errors.New("invalid request")

type Service interface {
	AuthenticateSession(ctx context.Context, r Request) (token.EncryptedToken, error)
}

type service struct {
	repository            aggregate_session.Repository
	sessionTokenEncryptor token.Encryptor[session_token.Token]
	ttl                   time.Duration
}

func (x service) AuthenticateSession(ctx context.Context, r Request) (token.EncryptedToken, error) {
	if err := r.Validate(); err != nil {
		return token.EncryptedToken{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	var encryptedToken token.EncryptedToken

	if err := x.repository.Replace(ctx, r.GetSessionID(), func(
		ctx context.Context,
		current *aggregate_session.Session,
	) (*aggregate_session.Session, error) {
		now := time.Now().UTC()
		if !current.IsActive(now) {
			return nil, aggregate_session.ErrInvalidCredentials
		}

		if !current.IsAnonymous() {
			return nil, aggregate_session.ErrAlreadyAuthenticated
		}

		id, err := model_session_id.GenerateSessionID()
		if err != nil {
			return nil, err
		}

		expiresAt := now.Add(x.ttl)
		open := aggregate_session_command.NewOpen(now, expiresAt)
		authenticate := aggregate_session_command.NewAuthenticate(r.GetUserID(), now, expiresAt)
		session := aggregate_session.NewSession(id)

		if err := session.Open(open); err != nil {
			return nil, err
		}

		if err := session.Authenticate(authenticate); err != nil {
			return nil, err
		}

		encrypted, err := x.sessionTokenEncryptor.Encrypt(session_token.NewToken(
			session.GetSessionID(),
			session.GetUpdatedAt(),
			session.GetExpiresAt(),
		))
		if err != nil {
			return nil, err
		}

		encryptedToken = encrypted

		return session, nil
	}); err != nil {
		return token.EncryptedToken{}, err
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
