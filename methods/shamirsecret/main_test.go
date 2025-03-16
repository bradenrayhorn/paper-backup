package shamirsecret

import (
	"testing"

	assert "github.com/bradenrayhorn/paper-backup/internal/testutils"
)

func TestQREndToEnd(t *testing.T) {
	toPrint, err := Encode("abc", "1234", 5, 3)
	assert.NoErr(t, err)

	res, err := DecodeFromQR(toPrint.QRShares, "1234")
	assert.NoErr(t, err)
	assert.Equal(t, "abc", res)
}

func TestTextEndToEnd(t *testing.T) {
	toPrint, err := Encode("abc", "1234", 5, 3)
	assert.NoErr(t, err)

	res, err := DecodeFromText(toPrint.TextShares, "1234")
	assert.NoErr(t, err)
	assert.Equal(t, "abc", res)
}
