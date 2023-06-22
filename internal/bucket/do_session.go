package bucket

import (
	"fmt"
	"os"
	"strings"

	"github.com/KTachibanaM/mear/internal/do"
	"github.com/KTachibanaM/mear/internal/s3"
)

func GetDigitalOceanSpacesCredentialsFromEnv() (string, string, error) {
	access_key_id, exists := os.LookupEnv("AWS_ACCESS_KEY_ID")
	if !exists {
		return "", "", fmt.Errorf("AWS_ACCESS_KEY_ID is not set")
	}
	if !strings.HasPrefix(access_key_id, "DO") {
		return "", "", fmt.Errorf("AWS_ACCESS_KEY_ID doesn't look like a DigitalOcean crdential")
	}
	secret_access_key, exists := os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	if !exists {
		return "", "", fmt.Errorf("AWS_SECRET_ACCESS_KEY is not set")
	}
	return access_key_id, secret_access_key, nil
}

func NewDigitalOceanSpacesS3Session(dc_picker do.DigitalOceanDataCenterPicker, access_key_id, secrete_access_key string) *s3.S3Session {
	return s3.NewS3Session(
		fmt.Sprintf("https://%v.digitaloceanspaces.com", dc_picker.Pick()),
		"us-east-1",
		access_key_id,
		secrete_access_key,
		false,
	)
}
