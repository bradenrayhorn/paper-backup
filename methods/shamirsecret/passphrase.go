package shamirsecret

import (
	"crypto/rand"

	"github.com/bradenrayhorn/paper-backup/encode"
)

func RandomPassphrase() (string, error) {
	random := make([]byte, 10)
	if _, err := rand.Read(random); err != nil {
		return "", err
	}

	return encode.ToSimpleText(random), nil
}
