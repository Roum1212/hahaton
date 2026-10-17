package service_get_or_register_volunteer

import (
	"context"
	"errors"
	"fmt"
	"time"

	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	aggregate_volunteer_command "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer/command"
	aggregate_wallet "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet"
	aggregate_wallet_command "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet/command"
	model_flida_profile "github.com/Roum1212/hahaton/internal/domain/model/flida-profile"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

var (
	ErrIncompleteProfile = errors.New("identity profile is incomplete")
	ErrInvalidIDToken    = errors.New("invalid id token")
	ErrInvalidRequest    = errors.New("invalid request")
)

type FlidaClient interface {
	VerifyIDToken(ctx context.Context, idToken string) (VerifiedIdentity, error)
}

type Service interface {
	GetOrRegisterVolunteer(ctx context.Context, r Request) (*aggregate_volunteer.Volunteer, error)
}

type service struct {
	flidaClient FlidaClient
	volunteers  aggregate_volunteer.Repository
	wallets     aggregate_wallet.Repository
}

func (x service) GetOrRegisterVolunteer(ctx context.Context, r Request) (*aggregate_volunteer.Volunteer, error) {
	if err := r.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	verified, err := x.flidaClient.VerifyIDToken(ctx, r.GetIDToken())
	if err != nil {
		return nil, err
	}

	volunteer, err := x.volunteers.FindByIdentity(ctx, verified.GetIdentity())
	switch {
	case err == nil:
		return volunteer, nil
	case !errors.Is(err, aggregate_volunteer.ErrNotFound):
		return nil, err
	}

	profile, ok := verified.GetProfile().Get()
	if !ok {
		return nil, ErrIncompleteProfile
	}

	volunteer, err = registerVolunteer(profile)
	if err != nil {
		return nil, err
	}

	wallet := aggregate_wallet.NewWallet(volunteer.GetUserID())
	if err := wallet.Open(aggregate_wallet_command.NewOpen(volunteer.GetCreatedAt())); err != nil {
		return nil, err
	}

	if err := x.volunteers.Create(ctx, volunteer, func(ctx context.Context) error {
		return x.wallets.Create(ctx, wallet)
	}); err != nil {
		return nil, err
	}

	return volunteer, nil
}

func registerVolunteer(profile model_flida_profile.Profile) (*aggregate_volunteer.Volunteer, error) {
	userID, err := model_user_id.GenerateUserID()
	if err != nil {
		return nil, err
	}

	command := aggregate_volunteer_command.NewRegister(
		profile.GetIdentity(), profile.GetFullName(), profile.GetDateOfBirth(),
		profile.GetEmailAddress(), profile.GetPhoneNumber(), time.Now().UTC(),
	)
	volunteer := aggregate_volunteer.NewVolunteer(userID)

	if err := volunteer.Register(command); err != nil {
		return nil, err
	}

	return volunteer, nil
}

func NewService(
	flidaClient FlidaClient,
	volunteers aggregate_volunteer.Repository,
	wallets aggregate_wallet.Repository,
) Service {
	return service{
		flidaClient: flidaClient,
		volunteers:  volunteers,
		wallets:     wallets,
	}
}
