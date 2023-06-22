package bucket

import (
	"github.com/KTachibanaM/mear/internal/s3"
)

var DevContainerS3Session = s3.NewS3Session(
	"http://minio:9000",
	"us-east-1",
	"minioadmin",
	"minioadmin",
	true,
)
