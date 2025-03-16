package crypt

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
)

type encryptionType = uint16

const (
	typeSingleFileAES256_V1 encryptionType = 0
	typePlaceholder_1       encryptionType = 1
	typePlaceholder_2       encryptionType = 2
	typePlaceholder_3       encryptionType = 3
	typePlaceholder_4       encryptionType = 4
	typePlaceholder_5       encryptionType = 5
)
const encryptTypes = 6 // count of above types

var errDecode = errors.New("invalid key or data")

func encodeHeader(passphrase string, kind encryptionType) []byte {
	sum := sha256.Sum256(append([]byte(passphrase), binary.BigEndian.AppendUint16([]byte{}, kind)...))
	return sum[:]
}

func decodeHeader(passphrase string, data []byte) (encryptionType, error) {
	if data == nil {
		return 0, errDecode
	}
	header := data[:sha256.Size]

	for i := 0; i < encryptTypes; i++ {
		kind := encryptionType(i)
		if bytes.Equal(header, encodeHeader(passphrase, kind)) {
			return kind, nil
		}
	}

	return 0, errDecode
}

func stripHeader(data []byte) []byte {
	if data == nil {
		return data
	}
	return data[sha256.Size:]
}
