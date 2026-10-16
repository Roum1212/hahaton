package paseto

import (
	"time"

	"aidanwoods.dev/go-paseto"
)

func ValidNow() paseto.Rule {
	return func(t paseto.Token) error {
		return paseto.ValidAt(time.Now().UTC())(t)
	}
}
