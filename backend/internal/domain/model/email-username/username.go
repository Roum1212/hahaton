package model_email_username

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	model_email_domain "github.com/Roum1212/hahaton/internal/domain/model/email-domain"
)

const maxUsernameLength = 64 // RFC 5321 section-4.5.3.1.1.

var protonReplacer = strings.NewReplacer(".", "", "-", "", "_", "")

type Username struct {
	string
}

func (x Username) IsZero() bool {
	return x.string == ""
}

func (x Username) Normalize(domain model_email_domain.Domain) Username {
	switch domain.Normalize().String() {
	case "gmail.com":
		return Username{strings.ToLower(strings.ReplaceAll(x.string, ".", ""))}
	case "yandex.ru":
		return Username{strings.ToLower(strings.ReplaceAll(x.string, "-", "."))}
	case "proton.me":
		return Username{strings.ToLower(protonReplacer.Replace(x.string))}
	default:
		return Username{strings.ToLower(x.string)}
	}
}

func (x Username) String() string {
	return x.string
}

func (x Username) Validate() error {
	if x.IsZero() {
		return errors.New("username is empty")
	}

	if len(x.String()) > maxUsernameLength {
		return errors.New("invalid username length")
	}

	if strings.HasPrefix(x.string, ".") || strings.HasPrefix(x.string, "-") {
		return errors.New("username starts with dot or hyphen")
	}

	if strings.HasSuffix(x.string, ".") || strings.HasSuffix(x.string, "-") {
		return errors.New("username ends with dot or hyphen")
	}

	if strings.Contains(x.string, "..") {
		return errors.New("username contains consecutive dots")
	}

	for _, r := range x.string {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			switch r {
			case '-', '.', '_':
			default:
				return fmt.Errorf("username contains invalid character: %q", r)
			}
		}
	}

	return nil
}

func NewUsername(s string) (Username, error) {
	x := Username{s}
	if err := x.Validate(); err != nil {
		return Username{}, fmt.Errorf("invalid username: %w", err)
	}

	return x, nil
}
