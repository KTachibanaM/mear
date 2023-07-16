package utils

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"
)

var DigitalOceanSpacesBucketNameMaxLength = 63
var DigitalOceanSpacesBucketSuffixLength = 8

func GetRandomName(prefix string, suffix_length, max_length int) (string, error) {
	if len(prefix) > DigitalOceanSpacesBucketNameMaxLength-DigitalOceanSpacesBucketSuffixLength-1 {
		return "", fmt.Errorf("prefix is too long")
	}
	suffix, err := randomSuffix(DigitalOceanSpacesBucketSuffixLength)
	if err != nil {
		return "", err
	}
	return prefix + "-" + suffix, nil
}

func randomSuffix(length int) (string, error) {
	// Generate a random byte slice
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Convert the byte slice to a base32-encoded string
	name := strings.Replace(strings.ToLower(base32.HexEncoding.EncodeToString(bytes)[:length]), "=", "0", -1)

	return name, nil
}
