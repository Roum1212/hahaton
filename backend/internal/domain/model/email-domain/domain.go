package model_email_domain

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"golang.org/x/net/publicsuffix"

	model_email_domain_label "github.com/Roum1212/hahaton/internal/domain/model/email-domain-label"
)

const maxDomainLength = 253 // RFC 5321 section-4.5.3.1.2.

var (
	gmail  = mustNewDomain("gmail.com")
	icloud = mustNewDomain("icloud.com")
	proton = mustNewDomain("proton.me")
	yandex = mustNewDomain("yandex.ru")
)

type Domain struct {
	labels []model_email_domain_label.Label
}

func (x Domain) Equals(v Domain) bool {
	if len(x.labels) != len(v.labels) {
		return false
	}

	for i := range x.labels {
		if x.labels[i] != v.labels[i] {
			return false
		}
	}

	return true
}

func (x Domain) GetLabels() []model_email_domain_label.Label {
	return slices.Clone(x.labels)
}

func (x Domain) IsZero() bool {
	return len(x.labels) == 0
}

func (x Domain) Normalize() Domain {
	switch x.String() {
	case "googlemail.com":
		return gmail
	case "mac.com", "me.com":
		return icloud
	case "narod.ru", "ya.ru", "yandex.by", "yandex.com", "yandex.com.am", "yandex.com.ge", "yandex.com.tr",
		"yandex.fr", "yandex.kz", "yandex.ua", "yandex.uz":
		return yandex
	case "pm.me", "protonmail.ch", "protonmail.com":
		return proton
	default:
		return x
	}
}

func (x Domain) String() string {
	labelsString := make([]string, len(x.labels))
	for i := range x.labels {
		labelsString[i] = x.labels[i].String()
	}

	return strings.Join(labelsString, ".")
}

func (x Domain) Validate() error {
	if x.IsZero() {
		return errors.New("domain is empty")
	}

	if len(x.String()) > maxDomainLength {
		return errors.New("invalid domain length")
	}

	if _, icann := publicsuffix.PublicSuffix(x.String()); !icann {
		return errors.New("invalid tld")
	}

	for _, label := range x.labels {
		if err := label.Validate(); err != nil {
			return fmt.Errorf("failed to validate label %q: %w", label, err)
		}
	}

	return nil
}

func mustNewDomain(s string) Domain {
	x, err := NewDomain(s)
	if err != nil {
		panic(fmt.Sprintf("invalid canonical domain %q: %v", s, err))
	}

	return x
}

func parseLabels(s string) ([]model_email_domain_label.Label, error) {
	labelsString := strings.Split(s, ".")
	if len(labelsString) < 2 { //nolint:mnd
		return nil, errors.New("must contain at least one dot")
	}

	labels := make([]model_email_domain_label.Label, len(labelsString))
	for i := range labelsString {
		label, err := model_email_domain_label.NewLabel(labelsString[i])
		if err != nil {
			return nil, fmt.Errorf("failed to create label: %w", err)
		}

		labels[i] = label
	}

	return labels, nil
}

func NewDomain(s string) (Domain, error) {
	labels, err := parseLabels(s)
	if err != nil {
		return Domain{}, fmt.Errorf("failed to parse domain labels: %w", err)
	}

	x := Domain{labels: labels}
	if err = x.Validate(); err != nil {
		return Domain{}, fmt.Errorf("invalid domain: %w", err)
	}

	return x, nil
}
