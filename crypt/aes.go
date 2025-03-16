package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const saltLength = 48

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

	encrypted := gcm.Seal(nil, nonce, data, nil)

	result := encodeHeader(passphrase, typeSingleFileAES256_V1)
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, encrypted...)

	return result, nil
}

func Decrypt(passphrase string, data []byte) ([]byte, error) {
	kind, err := decodeHeader(passphrase, data)
	if err != nil {
		return nil, err
	}

	if kind != typeSingleFileAES256_V1 {
		return nil, fmt.Errorf("cannot handle type %d", kind)
	}

	encrypted := stripHeader(data)
	if len(encrypted) < saltLength+1 {
		return nil, fmt.Errorf("crypto data invalid")
	}

	salt := encrypted[:saltLength]
	encrypted = encrypted[saltLength:]

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(encrypted) < gcm.NonceSize() {
		return nil, fmt.Errorf("crypto data invalid")
	}
	nonce := encrypted[:gcm.NonceSize()]
	encrypted = encrypted[gcm.NonceSize():]

	return gcm.Open(nil, nonce, encrypted, nil)
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
