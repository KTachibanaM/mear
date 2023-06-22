package bucket

import (
	"fmt"

	"github.com/KTachibanaM/mear/internal/do"
	"github.com/KTachibanaM/mear/internal/s3"
)

var DevContainerS3Session = s3.NewS3Session(
	"http://minio:9000",
	"us-east-1",
	"minioadmin",
	"minioadmin",
	true,
)

func NewDigitalOceanSpacesS3Session(dc_picker do.DigitalOceanDataCenterPicker, access_key_id, secrete_access_key string) *s3.S3Session {
	return s3.NewS3Session(
		fmt.Sprintf("https://%v.digitaloceanspaces.com", dc_picker.Pick()),
		"us-east-1",
		access_key_id,
		secrete_access_key,
		false,
	)
}
