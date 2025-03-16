package kind

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	mathRand "math/rand/v2"
)

type encryptionType = uint16

const (
	TypeSingleFile_V1   encryptionType = 0
	TypeShamirSecret_V1 encryptionType = 1
	TypePlaceholder_2   encryptionType = 2
	TypePlaceholder_3   encryptionType = 3
	TypePlaceholder_4   encryptionType = 4
	TypePlaceholder_5   encryptionType = 5
)

var errDecode = errors.New("unknown kind")

const headerSize = 4

func Encode(kind encryptionType) ([]byte, error) {
	seed := make([]byte, 2)
	_, err := rand.Read(seed)
	if err != nil {
		return nil, err
	}

	version := make([]byte, 2)
	binary.BigEndian.PutUint16(version, kind)
	mask(version, seed)

	return append(seed, version...), nil
}

func Decode(data []byte) (encryptionType, error) {
	if len(data) < headerSize {
		return 0, errDecode
	}

	seed := data[0:2]
	version := data[2:4]
	unmask(version, seed)

	return binary.BigEndian.Uint16(version), nil
}

func Strip(data []byte) []byte {
	if data == nil {
		return data
	}
	return data[headerSize:]
}

func makeRandomMap(seed []byte) []byte {
	chaChaSeed := [32]byte{}
	copy(chaChaSeed[:], seed)

	random := make([]byte, 256)
	for i := range random {
		random[i] = byte(i)
	}
	r := mathRand.New(mathRand.NewChaCha8(chaChaSeed))
	for i := len(random) - 1; i > 0; i-- {
		j := r.IntN(i + 1)
		random[i], random[j] = random[j], random[i]
	}
	return random
}

func mask(data, seed []byte) {
	random := makeRandomMap(seed)
	for i, b := range data {
		data[i] = random[b]
	}
}

func unmask(data, seed []byte) {
	random := makeRandomMap(seed)
	reverse := make([]byte, len(random))
	for i, v := range random {
		reverse[v] = byte(i)
	}
	for i, b := range data {
		data[i] = reverse[b]
	}
}
