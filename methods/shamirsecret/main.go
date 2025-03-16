package shamirsecret

import (
	"fmt"

	"github.com/bradenrayhorn/paper-backup/compress"
	"github.com/bradenrayhorn/paper-backup/crypt"
	"github.com/bradenrayhorn/paper-backup/encode"
	"github.com/bradenrayhorn/paper-backup/kind"
	"github.com/bradenrayhorn/paper-backup/shamir"
)

type encodeResult struct {
	QRShares   [][]byte
	TextShares []string
}

const LINE_SIZE = 40

func Encode(secret string, key string, parts int, threshold int) (encodeResult, error) {
	compressed, err := compress.Compress([]byte(secret))
	if err != nil {
		return encodeResult{}, fmt.Errorf("compressing: %w", err)
	}

	encrypted, err := crypt.Encrypt(key, compressed)
	if err != nil {
		return encodeResult{}, fmt.Errorf("encrypting: %w", err)
	}

	shares, err := shamir.Split(encrypted, parts, threshold)
	if err != nil {
		return encodeResult{}, fmt.Errorf("split: %w", err)
	}

	kind, err := kind.Encode(kind.TypeShamirSecret_V1)
	if err != nil {
		return encodeResult{}, fmt.Errorf("kind: %w", err)
	}

	qrShares := [][]byte{}
	textShares := []string{}

	for _, share := range shares {
		share = append(kind, share...)

		qrShares = append(qrShares, share)
		textShares = append(textShares, encode.ToText(share, LINE_SIZE))
	}

	return encodeResult{
		QRShares:   qrShares,
		TextShares: textShares,
	}, nil
}

func DecodeFromQR(shares [][]byte, key string) (string, error) {
	return decodeShares(shares, key)
}

func DecodeFromText(shares []string, key string) (string, error) {
	byteShares := make([][]byte, len(shares))

	for i, share := range shares {
		share, err := encode.FromText(share)
		if err != nil {
			return "", fmt.Errorf("decode: %w", err)
		}
		byteShares[i] = share
	}

	return decodeShares(byteShares, key)
}

func decodeShares(shares [][]byte, key string) (string, error) {
	strippedShares := make([][]byte, len(shares))
	for i, share := range shares {
		shareKind, err := kind.Decode(share)
		if err != nil {
			return "", fmt.Errorf("kind: %w", err)
		}

		if shareKind != kind.TypeShamirSecret_V1 {
			return "", fmt.Errorf("unknown kind: %d", shareKind)
		}

		strippedShares[i] = kind.Strip(share)
	}

	result, err := shamir.Combine(strippedShares)
	if err != nil {
		return "", fmt.Errorf("combine: %w", err)
	}

	decrypted, err := crypt.Decrypt(key, result)
	if err != nil {
		return "", fmt.Errorf("decrypt: %w", err)
	}

	decompressed, err := compress.Decompress(decrypted)
	if err != nil {
		return "", fmt.Errorf("decompress: %w", err)
	}

	return string(decompressed), nil
}
