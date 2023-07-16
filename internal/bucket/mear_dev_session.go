package bucket

import (
	"github.com/KTachibanaM/mear/internal/s3"
)

var MearDevSession = s3.NewS3Session(
	"",
	"us-west-2",
	"",
	"",
	false,
)
