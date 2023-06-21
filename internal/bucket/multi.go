package bucket

import (
	"github.com/KTachibanaM/mear/internal/s3"
	"github.com/hashicorp/go-multierror"
)

type MultiBucketProvisioner struct {
	provisioned []*s3.S3Bucket
}

func NewMultiBucketProvisioner() *MultiBucketProvisioner {
	return &MultiBucketProvisioner{
		provisioned: []*s3.S3Bucket{},
	}
}

func (p *MultiBucketProvisioner) Provision(buckets []*s3.S3Bucket) error {
	var errors error
	for _, bucket := range buckets {
		err := ProvisionBucket(bucket)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			p.provisioned = append(p.provisioned, bucket)
		}
	}
	return errors
}

func (p *MultiBucketProvisioner) Teardown() error {
	var errors error
	for _, bucket := range p.provisioned {
		if err := TeardownBucket(bucket); err != nil {
			errors = multierror.Append(errors, err)
		}
	}
	return errors
}
