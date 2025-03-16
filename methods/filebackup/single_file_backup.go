package filebackup

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/bradenrayhorn/paper-backup/compress"
	"github.com/bradenrayhorn/paper-backup/crypt"
	"github.com/bradenrayhorn/paper-backup/kind"
)

type fileBackup struct {
	Name string
	Data []byte
}

func Encode(data []byte, fileName string, key string) ([]byte, error) {
	compressed, err := compress.Compress(data)
	if err != nil {
		return nil, fmt.Errorf("compressing: %w", err)
	}

	var encoded bytes.Buffer
	backup := fileBackup{Name: fileName, Data: compressed}
	enc := gob.NewEncoder(&encoded)
	if err := enc.Encode(backup); err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}

	encrypted, err := crypt.Encrypt(key, encoded.Bytes())
	if err != nil {
		return nil, fmt.Errorf("encrypting: %w", err)
	}

	kind, err := kind.Encode(kind.TypeSingleFile_V1)
	if err != nil {
		return nil, fmt.Errorf("header: %w", err)
	}

	return append(kind, encrypted...), nil
}

func Decode(data []byte, key string) ([]byte, string, error) {
	fileKind, err := kind.Decode(data)
	if err != nil {
		return nil, "", fmt.Errorf("check kind: %w", err)
	}

	if fileKind != kind.TypeSingleFile_V1 {
		return nil, "", fmt.Errorf("unexpected kind: %d", fileKind)
	}

	decrypted, err := crypt.Decrypt(key, kind.Strip(data))
	if err != nil {
		return nil, "", fmt.Errorf("decrypt: %w", err)
	}

	dec := gob.NewDecoder(bytes.NewReader(decrypted))
	var backup fileBackup
	if err := dec.Decode(&backup); err != nil {
		return nil, "", fmt.Errorf("decode: %w", err)
	}

	result, err := compress.Decompress(backup.Data)
	if err != nil {
		return nil, "", fmt.Errorf("decompress: %w", err)
	}

	return result, backup.Name, nil
}
