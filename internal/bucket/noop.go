package bucket

import (
	"fmt"

	"github.com/KTachibanaM/mear/internal/s3"
	"github.com/aws/aws-sdk-go/aws"

	aws_s3 "github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
)

func deleteS3Target(s3_target *s3.S3Target) error {
	s3_sess, err := s3.CreateS3Session(s3_target)
	if err != nil {
		return err
	}
	_, err = aws_s3.New(s3_sess).DeleteObject(
		&aws_s3.DeleteObjectInput{
			Bucket: aws.String(s3_target.BucketName),
			Key:    aws.String(s3_target.ObjectKey),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

type NoOpBucketProvisioner struct {
	s3_target           *s3.S3Target
	delete_at_provision bool
}

func NewNoOpBucketProvisioner(s3_target *s3.S3Target, delete_at_provision bool) *NoOpBucketProvisioner {
	return &NoOpBucketProvisioner{
		s3_target:           s3_target,
		delete_at_provision: delete_at_provision,
	}
}

func (p *NoOpBucketProvisioner) Provision() (*s3.S3Target, error) {
	if p.delete_at_provision {
		log.Infof("at provision deleting s3 target %v", p.s3_target)
		if err := deleteS3Target(p.s3_target); err != nil {
			return nil, fmt.Errorf("at provision failed to delete s3 target %v: %v", p.s3_target, err)
		}
	}
	return p.s3_target, nil
}

func (p *NoOpBucketProvisioner) Teardown() error {
	log.Infof("fake deleting s3 target %v", p.s3_target)
	return nil
}
