package model_email_domain_label

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/net/idna"
)

const maxLabelLength = 63

type Label struct {
	string
}

func (x Label) IsZero() bool {
	return x.string == ""
}

func (x Label) String() string {
	return x.string
}

func (x Label) Validate() error {
	if x.IsZero() {
		return errors.New("label is empty")
	}

	if len(x.string) > maxLabelLength {
		return errors.New("invalid label length")
	}

	for _, r := range x.string {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' { // RFC 1035, RFC 1123.
			return fmt.Errorf("label contains invalid character: %q", r)
		}
	}

	return nil
}

func NewLabel(s string) (Label, error) {
	s, err := idna.Registration.ToASCII(strings.ToLower(s))
	if err != nil {
		return Label{}, fmt.Errorf("failed to convert to ascii: %w", err)
	}

	x := Label{s}
	if err = x.Validate(); err != nil {
		return Label{}, fmt.Errorf("invalid label: %w", err)
	}

	return x, nil
}
