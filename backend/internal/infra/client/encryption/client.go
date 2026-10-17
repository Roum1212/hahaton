package client_encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

type Client struct {
	aead cipher.AEAD
}

func (x Client) Decrypt(value []byte) (string, error) {
	size := x.aead.NonceSize()
	if len(value) < size {
		return "", errors.New("invalid encrypted value")
	}

	text, err := x.aead.Open(nil, value[:size], value[size:], nil)
	if err != nil {
		return "", err
	}

	return string(text), nil
}

func (x Client) Encrypt(text string) ([]byte, error) {
	nonce := make([]byte, x.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return x.aead.Seal(nonce, nonce, []byte(text), nil), nil
}

func NewClient(key []byte) (Client, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return Client{}, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return Client{}, err
	}

	return Client{aead: aead}, nil
}
