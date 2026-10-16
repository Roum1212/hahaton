package token

import (
	"errors"
	"fmt"
)

type Service[T any] interface {
	Decryptor[T]
	Encryptor[T]
}

type Decryptor[T any] interface {
	Decrypt(encryptedToken EncryptedToken) (T, error)
}

type Encryptor[T any] interface {
	Encrypt(value T) (EncryptedToken, error)
}

type EncryptedToken struct {
	value string
}

func (x EncryptedToken) String() string {
	return x.value
}

func (x EncryptedToken) Validate() error {
	if len(x.value) == 0 {
		return errors.New("encrypted token is empty")
	}

	return nil
}

func NewEncryptedToken(value string) (EncryptedToken, error) {
	x := EncryptedToken{value: value}
	if err := x.Validate(); err != nil {
		return EncryptedToken{}, fmt.Errorf("invalid encrypted token: %w", err)
	}

	return x, nil
}
