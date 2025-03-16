package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const saltLength = 24

func Encrypt(passphrase string, data []byte) ([]byte, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	cipher := gcm.Seal(nil, nonce, data, nil)

	result := []byte{}
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, cipher...)

	return result, nil
}

func Decrypt(passphrase string, data []byte) ([]byte, error) {
	salt := data[:saltLength]
	cipherWithNonce := data[saltLength:]

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(cipherWithNonce) < gcm.NonceSize() {
		return nil, fmt.Errorf("crypto data invalid")
	}
	nonce := cipherWithNonce[:gcm.NonceSize()]
	cipher := cipherWithNonce[gcm.NonceSize():]

	return gcm.Open(nil, nonce, cipher, nil)
}

func deriveKey(passphrase string, salt []byte) []byte {
	return argon2.IDKey(
		[]byte(passphrase),
		salt,
		8,
		512*1024,
		8,
		32,
	)
}
