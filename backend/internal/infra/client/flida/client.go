package client_flida

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"

	service_get_or_register_volunteer "github.com/Roum1212/hahaton/internal/app/service/get-or-register-volunteer"
	model_date_of_birth "github.com/Roum1212/hahaton/internal/domain/model/date-of-birth"
	model_email_address "github.com/Roum1212/hahaton/internal/domain/model/email-address"
	model_flida_profile "github.com/Roum1212/hahaton/internal/domain/model/flida-profile"
	model_full_name "github.com/Roum1212/hahaton/internal/domain/model/full-name"
	model_identity "github.com/Roum1212/hahaton/internal/domain/model/identity"
	model_phone_number "github.com/Roum1212/hahaton/internal/domain/model/phone-number"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

type Client struct {
	verifier *oidc.IDTokenVerifier
	clientID string
}

func (x Client) VerifyIDToken(
	ctx context.Context,
	idToken string,
) (service_get_or_register_volunteer.VerifiedIdentity, error) {
	verified, err := x.verifier.Verify(ctx, idToken)
	if err != nil {
		return service_get_or_register_volunteer.VerifiedIdentity{}, fmt.Errorf(
			"%w: %w",
			service_get_or_register_volunteer.ErrInvalidIDToken,
			err,
		)
	}

	if verified.IssuedAt.IsZero() || verified.IssuedAt.After(time.Now().UTC()) {
		return service_get_or_register_volunteer.VerifiedIdentity{}, service_get_or_register_volunteer.ErrInvalidIDToken
	}

	var claims map[string]any

	if err := verified.Claims(&claims); err != nil {
		return service_get_or_register_volunteer.VerifiedIdentity{}, service_get_or_register_volunteer.ErrInvalidIDToken
	}

	authorizedParty := claimString(claims, "azp")
	if (len(verified.Audience) > 1 || len(authorizedParty) != 0) && authorizedParty != x.clientID {
		return service_get_or_register_volunteer.VerifiedIdentity{}, service_get_or_register_volunteer.ErrInvalidIDToken
	}

	identity, err := model_identity.NewIdentity(verified.Issuer, verified.Subject)
	if err != nil {
		return service_get_or_register_volunteer.VerifiedIdentity{}, service_get_or_register_volunteer.ErrInvalidIDToken
	}

	var profile utils_value.Optional[model_flida_profile.Profile]

	// Existing volunteers can sign in without repeating registration claims.
	if value, err := profileFromClaims(identity, claims); err == nil {
		profile = utils_value.NewOptional(value)
	}

	return service_get_or_register_volunteer.NewVerifiedIdentity(identity, profile), nil
}

func claimString(claims map[string]any, key string) string {
	value, ok := claims[key].(string)
	if !ok {
		return ""
	}

	return value
}

func profileFromClaims(identity model_identity.Identity, claims map[string]any) (model_flida_profile.Profile, error) {
	fullName, err := model_full_name.NewFullName(
		claimString(claims, "given_name"), claimString(claims, "family_name"), claimString(claims, "middle_name"),
	)
	if err != nil {
		return model_flida_profile.Profile{}, err
	}

	dateOfBirth, err := model_date_of_birth.NewDateOfBirthFromString(claimString(claims, "date_of_birth"))
	if err != nil {
		return model_flida_profile.Profile{}, err
	}

	emailAddress, err := model_email_address.NewEmailAddress(claimString(claims, "email"))
	if err != nil {
		return model_flida_profile.Profile{}, err
	}

	phoneNumber, err := model_phone_number.NewPhoneNumber(claimString(claims, "phone_number"))
	if err != nil {
		return model_flida_profile.Profile{}, err
	}

	return model_flida_profile.NewProfile(identity, fullName, dateOfBirth, emailAddress, phoneNumber)
}

func NewClient(ctx context.Context, issuer, clientID string, httpClient *http.Client) (Client, error) {
	issuerURL, err := url.Parse(issuer)
	if err != nil || issuerURL.Scheme != "https" || len(issuerURL.Host) == 0 {
		return Client{}, errors.New("flida issuer must be an absolute https url")
	}

	if len(strings.TrimSpace(clientID)) == 0 {
		return Client{}, errors.New("flida client id is empty")
	}

	if httpClient == nil || httpClient.Timeout <= 0 {
		return Client{}, errors.New("flida http client must have a positive timeout")
	}

	provider, err := oidc.NewProvider(oidc.ClientContext(ctx, httpClient), issuer)
	if err != nil {
		return Client{}, fmt.Errorf("discover flida provider: %w", err)
	}

	return Client{
		verifier: provider.Verifier(&oidc.Config{
			ClientID:             clientID,
			SupportedSigningAlgs: []string{oidc.RS256},
		}),
		clientID: clientID,
	}, nil
}
