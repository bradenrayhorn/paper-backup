package crypt

import (
	"testing"

	assert "github.com/bradenrayhorn/paper-backup/internal/testutils"
)

func TestDifferentiatesKinds(t *testing.T) {
	header := encodeHeader("abc", typeSingleFileAES256_V1)

	result, err := decodeHeader("abc", header)
	assert.NoErr(t, err)
	assert.Equal(t, typeSingleFileAES256_V1, result)

	header = encodeHeader("abc", typePlaceholder_3)
	result, err = decodeHeader("abc", header)
	assert.NoErr(t, err)
	assert.Equal(t, typePlaceholder_3, result)
}

func TestFailsOnNonMatching(t *testing.T) {
	header := encodeHeader("abc", typeSingleFileAES256_V1)

	_, err := decodeHeader("abcd", header)
	assert.ErrContains(t, err, "invalid key or data")
}
