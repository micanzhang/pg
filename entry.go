package pg

import (
	"crypto/sha512"
	"encoding/base64"
)

type Password string

func NewPassword(key string, plaintext string) (Password, error) {
	data, err := Encrypt(Key(key), []byte(plaintext))
	if err != nil {
		return "", err
	}

	return Password(base64.StdEncoding.EncodeToString(data)), nil
}

func (p Password) Plaintext(key string) (string, error) {
	bs, err := base64.StdEncoding.DecodeString(string(p))
	if err != nil {
		return "", err
	}

	data, err := Decrypt(Key(key), bs)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

type Entry struct {
	Domain   string   `json:"domain"`
	Username string   `json:"username"`
	Password Password `json:"password"`
}

func Key(key string) []byte {
	dst := sha512.Sum512([]byte(key))
	return dst[:]
}
