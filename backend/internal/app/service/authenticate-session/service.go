package service_authenticate_session

import (
	"context"
	"errors"
	"fmt"
	"time"

	aggregate_session "github.com/Roum1212/hahaton/internal/domain/aggregate/session"
	aggregate_session_command "github.com/Roum1212/hahaton/internal/domain/aggregate/session/command"
	model_session_id "github.com/Roum1212/hahaton/internal/domain/model/session-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	session_token "github.com/Roum1212/hahaton/internal/pkg/session/token"
	"github.com/Roum1212/hahaton/internal/pkg/token"
)

var (
	ErrIncompleteProfile = errors.New("identity profile is incomplete")
	ErrInvalidIDToken    = errors.New("invalid id token")
	ErrInvalidRequest    = errors.New("invalid request")
)

type Authenticator interface {
	Authenticate(ctx context.Context, idToken string) (model_user_id.UserID, error)
}

type Transactor interface {
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type Service interface {
	AuthenticateSession(ctx context.Context, r Request) (token.EncryptedToken, error)
}

type service struct {
	repository            aggregate_session.Repository
	authenticator         Authenticator
	transactor            Transactor
	sessionTokenEncryptor token.Encryptor[session_token.Token]
	ttl                   time.Duration
}

func (x service) AuthenticateSession(ctx context.Context, r Request) (token.EncryptedToken, error) {
	if err := r.Validate(); err != nil {
		return token.EncryptedToken{}, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	sessionID := r.GetSessionID()

	current, err := x.repository.Load(ctx, sessionID)
	if err != nil {
		return token.EncryptedToken{}, err
	}

	if !current.IsActive(time.Now().UTC()) {
		return token.EncryptedToken{}, aggregate_session.ErrInvalidCredentials
	}

	if !current.GetUserID().IsZero() {
		return token.EncryptedToken{}, aggregate_session.ErrAlreadyAuthenticated
	}

	userID, err := x.authenticator.Authenticate(ctx, r.GetIDToken())
	if err != nil {
		return token.EncryptedToken{}, err
	}

	id, err := model_session_id.GenerateSessionID()
	if err != nil {
		return token.EncryptedToken{}, err
	}

	now := time.Now().UTC()
	expiresAt := now.Add(x.ttl)

	open := aggregate_session_command.NewOpen(now, expiresAt)
	authenticate := aggregate_session_command.NewAuthenticate(userID, now, expiresAt)
	session := aggregate_session.NewSession(id)

	if err = session.Open(open); err != nil {
		return token.EncryptedToken{}, err
	}

	if err = session.Authenticate(authenticate); err != nil {
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

	if err = x.transactor.InTx(ctx, func(ctx context.Context) error {
		if deleteErr := x.repository.Delete(ctx, sessionID); deleteErr != nil {
			return deleteErr
		}

		return x.repository.Create(ctx, session)
	}); err != nil {
		return token.EncryptedToken{}, err
	}

	return encryptedToken, nil
}

func NewService(
	repository aggregate_session.Repository,
	authenticator Authenticator,
	transactor Transactor,
	sessionTokenEncryptor token.Encryptor[session_token.Token],
	ttl time.Duration,
) Service {
	return service{
		repository:            repository,
		authenticator:         authenticator,
		transactor:            transactor,
		sessionTokenEncryptor: sessionTokenEncryptor,
		ttl:                   ttl,
	}
}
