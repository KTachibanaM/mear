package bucket

import (
	"github.com/KTachibanaM/mear/internal/s3"
	"github.com/hashicorp/go-multierror"
)

type MultiBucketProvisioner struct {
	provisioned_buckets []*s3.S3Bucket
}

func NewMultiBucketProvisioner() *MultiBucketProvisioner {
	return &MultiBucketProvisioner{
		provisioned_buckets: []*s3.S3Bucket{},
	}
}

func (p *MultiBucketProvisioner) Provision(buckets []*s3.S3Bucket) error {
	var errors error
	for _, bucket := range buckets {
		provisioned_bucket, err := NewBucketProvisioner(bucket.S3Session).Provision(bucket.BucketName)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			p.provisioned_buckets = append(p.provisioned_buckets, provisioned_bucket)
		}
	}
	return errors
}

func (p *MultiBucketProvisioner) Teardown() error {
	var errors error
	for _, bucket := range p.provisioned_buckets {
		if err := NewBucketProvisioner(bucket.S3Session).Teardown(); err != nil {
			errors = multierror.Append(errors, err)
		}
	}
	return errors
}
