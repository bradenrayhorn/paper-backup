package methods

import (
	"testing"

	assert "github.com/bradenrayhorn/paper-backup/internal/testutils"
)

func TestSingleFileBackup(t *testing.T) {
	data := []byte("abc")

	res, err := FileBackupEncode(data, "myfile.txt", "shh")
	assert.NoErr(t, err)

	recovered, name, err := FileBackupDecode(res, "shh")
	assert.NoErr(t, err)

	assert.Equal(t, "myfile.txt", name)
	assert.Equal(t, string(data), string(recovered))
}
