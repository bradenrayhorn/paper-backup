package main

import (
	"fmt"

	"github.com/bradenrayhorn/paper-backup/compress"
	"github.com/bradenrayhorn/paper-backup/crypt"
	"github.com/bradenrayhorn/paper-backup/encode"
)

func EncodeBackup(data []byte, key string) (string, []byte, error) {
	compressed, err := compress.Compress(data)
	if err != nil {
		return "", nil, fmt.Errorf("compressing: %w", err)
	}

	encrypted, err := crypt.Encrypt(key, compressed)
	if err != nil {
		return "", nil, fmt.Errorf("encrypting: %w", err)
	}

	return encode.ToText(encrypted), encrypted, nil
}

func DecodeBackupFromQR(data []byte, key string) ([]byte, error) {
	decrypted, err := crypt.Decrypt(key, data)
	if err != nil {
		return nil, err
	}

	return compress.Decompress(decrypted)
}

func DecodeBackupFromText(data string, key string) ([]byte, error) {
	decoded, err := encode.FromText(data)
	if err != nil {
		return nil, err
	}

	decrypted, err := crypt.Decrypt(key, decoded)
	if err != nil {
		return nil, err
	}

	return compress.Decompress(decrypted)
}
