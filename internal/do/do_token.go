package do

import (
	"fmt"
	"os"
)

func GetDigitalOceanTokenFromEnv() (string, error) {
	do_token, exists := os.LookupEnv("DO_TOKEN")
	if !exists {
		return "", fmt.Errorf("DO_TOKEN is not set")
	}
	return do_token, nil
}
