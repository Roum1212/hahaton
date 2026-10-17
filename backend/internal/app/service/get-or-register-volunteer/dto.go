package service_get_or_register_volunteer

import (
	"errors"
	"strings"

	model_flida_profile "github.com/Roum1212/hahaton/internal/domain/model/flida-profile"
	model_identity "github.com/Roum1212/hahaton/internal/domain/model/identity"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

type Request struct {
	idToken string
}

func (x Request) GetIDToken() string {
	return x.idToken
}

func (x Request) Validate() error {
	if len(strings.TrimSpace(x.idToken)) == 0 {
		return errors.New("id token is required")
	}

	return nil
}

func NewRequest(idToken string) Request {
	return Request{idToken: idToken}
}

type VerifiedIdentity struct {
	identity model_identity.Identity
	profile  utils_value.Optional[model_flida_profile.Profile]
}

func (x VerifiedIdentity) GetIdentity() model_identity.Identity {
	return x.identity
}

func (x VerifiedIdentity) GetProfile() utils_value.Optional[model_flida_profile.Profile] {
	return x.profile
}

func NewVerifiedIdentity(
	identity model_identity.Identity,
	profile utils_value.Optional[model_flida_profile.Profile],
) VerifiedIdentity {
	return VerifiedIdentity{
		identity: identity,
		profile:  profile,
	}
}
