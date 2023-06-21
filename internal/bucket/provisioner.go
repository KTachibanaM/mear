package bucket

import "github.com/KTachibanaM/mear/internal/s3"

type BucketProvisioner interface {
	Provision() (*s3.S3Target, error)
	Teardown() error
}
