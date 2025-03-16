package kind

import (
	"fmt"
	"testing"

	assert "github.com/bradenrayhorn/paper-backup/internal/testutils"
)

func TestCanEncodeAndDecodeHeader(t *testing.T) {
	header, err := Encode(1)
	assert.NoErr(t, err)
	fmt.Println(header)

	result, err := Decode(header)
	assert.NoErr(t, err)
	assert.Equal(t, 1, result)
}

func TestCanDecodeHeader(t *testing.T) {
	result, err := Decode([]byte{27, 177, 73, 196})
	assert.NoErr(t, err)
	assert.Equal(t, 1, result)
}

func TestFailsOnTooShort(t *testing.T) {
	_, err := Decode([]byte{1})
	assert.ErrIs(t, err, errDecode)
}

func TestFailsOnNil(t *testing.T) {
	_, err := Decode(nil)
	assert.ErrIs(t, err, errDecode)
}
