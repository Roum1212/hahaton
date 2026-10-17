package model_email_address

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"

	model_email_domain "github.com/Roum1212/hahaton/internal/domain/model/email-domain"
	model_email_username "github.com/Roum1212/hahaton/internal/domain/model/email-username"
)

const maxEmailAddressLength = 254 // RFC 5321 section-4.5.3.1.3 (без символов '<', '>').

type EmailAddress struct {
	domain   model_email_domain.Domain
	username model_email_username.Username
}

func (x EmailAddress) Equals(v EmailAddress) bool {
	return x.username == v.username && x.domain.Equals(v.domain)
}

func (x EmailAddress) GetDomain() model_email_domain.Domain {
	return x.domain
}

func (x EmailAddress) GetUsername() model_email_username.Username {
	return x.username
}

func (x EmailAddress) IsZero() bool {
	return x.domain.IsZero() && x.username.IsZero()
}

func (x EmailAddress) String() string {
	if x.IsZero() {
		return ""
	}

	return x.username.String() + "@" + x.domain.String()
}

func (x EmailAddress) Validate() error {
	if x.IsZero() {
		return errors.New("e-mail address is empty")
	}

	if len(x.String()) > maxEmailAddressLength {
		return errors.New("invalid e-mail address length")
	}

	if err := x.domain.Validate(); err != nil {
		return fmt.Errorf("failed to validate domain: %w", err)
	}

	if err := x.username.Validate(); err != nil {
		return fmt.Errorf("failed to validate username: %w", err)
	}

	return nil
}

func NewEmailAddress(s string) (EmailAddress, error) {
	mailAddress, err := mail.ParseAddress(s)
	if err != nil {
		return EmailAddress{}, fmt.Errorf("invalid e-mail address: %w", err)
	}

	localPartString, domainString, ok := strings.Cut(mailAddress.Address, "@")
	if !ok {
		return EmailAddress{}, errors.New("invalid e-mail address")
	}

	domain, err := model_email_domain.NewDomain(domainString)
	if err != nil {
		return EmailAddress{}, fmt.Errorf("failed to create domain: %w", err)
	}

	username, err := model_email_username.NewUsername(localPartString)
	if err != nil {
		return EmailAddress{}, fmt.Errorf("failed to create username: %w", err)
	}

	domain = domain.Normalize()
	username = username.Normalize(domain)

	x := EmailAddress{
		domain:   domain,
		username: username,
	}
	if err = x.Validate(); err != nil {
		return EmailAddress{}, fmt.Errorf("invalid e-mail address: %w", err)
	}

	return x, nil
}
