package main

import (
	"testing"

	assert "github.com/bradenrayhorn/paper-backup/internal/testutils"
)

func TestBackupFileToQR(t *testing.T) {
	data := []byte("abc")

	_, res, err := EncodeBackup(data, "shh")
	assert.NoErr(t, err)

	recovered, err := DecodeBackupFromQR(res, "shh")
	assert.NoErr(t, err)

	assert.Equal(t, string(data), string(recovered))
}
