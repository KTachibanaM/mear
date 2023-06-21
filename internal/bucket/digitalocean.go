package bucket

import (
	"fmt"

	"github.com/KTachibanaM/mear/internal/do"
	"github.com/KTachibanaM/mear/internal/s3"
	aws_s3 "github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-sdk-go/aws"
)

type DigitalOceanBucketProvisioner struct {
	dc_guesser        do.DigitalOceanDataCenterGuesser
	bucket_name       string
	object_key        string
	access_key_id     string
	secret_access_key string
}

func NewDigitalOceanBucketProvisioner(
	dc_guesser do.DigitalOceanDataCenterGuesser,
	bucket_name, object_key, access_key_id, secret_access_key string,
) *DigitalOceanBucketProvisioner {
	return &DigitalOceanBucketProvisioner{
		dc_guesser:        dc_guesser,
		bucket_name:       bucket_name,
		object_key:        object_key,
		access_key_id:     access_key_id,
		secret_access_key: secret_access_key,
	}
}

func (p *DigitalOceanBucketProvisioner) Provision() (*s3.S3Target, error) {
	s3_target := s3.NewS3Target(
		fmt.Sprintf("https://%v.digitaloceanspaces.com", p.dc_guesser.Guess()),
		"us-east-1",
		p.bucket_name,
		p.object_key,
		p.access_key_id,
		p.secret_access_key,
		false,
	)

	sess, err := s3.CreateS3Session(s3_target)
	if err != nil {
		return nil, fmt.Errorf("could not create S3 session for bucket provisioning: %w", err)
	}

	_, err = aws_s3.New(sess).CreateBucket(
		&aws_s3.CreateBucketInput{
			Bucket: aws.String(s3_target.BucketName),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("could not create bucket: %v", err)
	}
	return s3_target, nil

}

func (p *DigitalOceanBucketProvisioner) Teardown() error {
	s3_target := s3.NewS3Target(
		fmt.Sprintf("https://%v.digitaloceanspaces.com", p.dc_guesser.Guess()),
		"us-east-1",
		p.bucket_name,
		p.object_key,
		p.access_key_id,
		p.secret_access_key,
		false,
	)

	sess, err := s3.CreateS3Session(s3_target)
	if err != nil {
		return fmt.Errorf("could not create S3 session for bucket teardown: %w", err)
	}

	_, err = aws_s3.New(sess).DeleteBucket(
		&aws_s3.DeleteBucketInput{
			Bucket: aws.String(s3_target.BucketName),
		},
	)
	if err != nil {
		return fmt.Errorf("could not delete bucket: %v", err)
	}

	return nil
}
