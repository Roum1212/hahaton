package model_identity

import (
	"errors"
	"fmt"
	"strings"
)

type Identity struct {
	issuer  string
	subject string
}

func (x Identity) GetIssuer() string {
	return x.issuer
}

func (x Identity) GetSubject() string {
	return x.subject
}

func (x Identity) Validate() error {
	if len(strings.TrimSpace(x.issuer)) == 0 {
		return errors.New("issuer is empty")
	}

	if len(strings.TrimSpace(x.subject)) == 0 {
		return errors.New("subject is empty")
	}

	return nil
}

func NewIdentity(issuer, subject string) (Identity, error) {
	x := Identity{
		issuer:  issuer,
		subject: subject,
	}
	if err := x.Validate(); err != nil {
		return Identity{}, fmt.Errorf("invalid identity: %w", err)
	}

	return x, nil
}
