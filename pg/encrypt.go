package pg

import (
	"crypto/aes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"crypto/cipher"
)

const (
	NonceSize = aes.BlockSize
	MACSize   = 32 // Output size of HMAC-SHA-256
	CKeySize  = 32 // Cipher key size - AES-256
	MKeySize  = 32 // HMAC key size - HMAC-SHA-256
)

var KeySize = CKeySize + MKeySize

var (
	ErrEncrypt = errors.New("secret: encryption failed")
	ErrDecrypt = errors.New("secret: decryption failed")
)

func pad(in []byte) []byte {
	padding := aes.BlockSize - (len(in) % aes.BlockSize)
	for i := 0; i < padding; i++ {
		in = append(in, byte(padding))
	}

	return in
}

func unpad(in []byte) []byte {
	if len(in) == 0 {
		return nil
	}

	padding := in[len(in)-1]
	if int(padding) > len(in) || padding > aes.BlockSize {
		return nil
	} else if padding == 0 {
		return nil
	}

	for i := len(in) - 1; i > len(in)-int(padding)-1; i-- {
		if in[i] != padding {
			return nil
		}
	}

	return in[:len(in)-int(padding)]
}

func randBytes(n int) ([]byte, error) {
	r := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func Encrypt(key, message []byte) ([]byte, error) {
	if len(key) != KeySize {
		return nil, ErrEncrypt
	}

	iv, err := randBytes(NonceSize)
	if err != nil {
		return nil, ErrEncrypt
	}

	pmessage := pad(message)
	ct := make([]byte, len(pmessage))

	c, _ := aes.NewCipher(key[:CKeySize])
	ctr := cipher.NewCBCEncrypter(c, iv)
	ctr.CryptBlocks(ct, pmessage)

	h := hmac.New(sha256.New, key[CKeySize:])
	ct = append(iv, ct...)
	h.Write(ct)
	ct = h.Sum(ct)

	return ct, nil
}

func Decrypt(key, message []byte) ([]byte, error) {
	if len(key) != KeySize {
		return nil, ErrDecrypt
	}

	if (len(message) % aes.BlockSize) != 0 {
		return nil, ErrDecrypt
	}

	if len(message) < (4 * aes.BlockSize) {
		return nil, ErrDecrypt
	}

	macStart := len(message) - MACSize
	tag := message[macStart:]
	out := make([]byte, macStart-NonceSize)
	message = message[:macStart]

	h := hmac.New(sha256.New, key[CKeySize:])
	h.Write(message)
	mac := h.Sum(nil)
	if !hmac.Equal(mac, tag) {
		return nil, ErrDecrypt
	}

	c, _ := aes.NewCipher(key[:CKeySize])
	ctr := cipher.NewCBCDecrypter(c, message[:NonceSize])
	ctr.CryptBlocks(out, message[NonceSize:])

	pt := unpad(out)
	if pt == nil {
		return nil, ErrDecrypt
	}

	return pt, nil

}
