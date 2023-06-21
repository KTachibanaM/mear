package do

import (
	"crypto/rand"
	"encoding/base32"
	"strings"
)

func RandomBucketSuffix(length int) (string, error) {
	// Generate a random byte slice
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Convert the byte slice to a base32-encoded string
	name := strings.Replace(strings.ToLower(base32.HexEncoding.EncodeToString(bytes)[:length]), "=", "0", -1)

	return name, nil
}
