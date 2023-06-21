package bucket

import (
	"fmt"

	log "github.com/sirupsen/logrus"

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

	log.Println("creating bucket on DO...")
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

	log.Println("emptying bucket on DO...")
	resp, err := aws_s3.New(sess).ListObjectsV2(&aws_s3.ListObjectsV2Input{
		Bucket: aws.String(s3_target.BucketName),
	})
	if err != nil {
		return fmt.Errorf("could not list objects in bucket for deletion: %v", err)
	}
	var keys []*aws_s3.ObjectIdentifier
	for _, obj := range resp.Contents {
		keys = append(keys, &aws_s3.ObjectIdentifier{
			Key: obj.Key,
		})
	}
	if len(keys) > 0 {
		_, err = aws_s3.New(sess).DeleteObjects(&aws_s3.DeleteObjectsInput{
			Bucket: aws.String(s3_target.BucketName),
			Delete: &aws_s3.Delete{
				Objects: keys,
				Quiet:   aws.Bool(true),
			},
		})
		if err != nil {
			return fmt.Errorf("could not delete objects in bucket: %v", err)
		}
	}

	log.Println("deleting bucket on DO...")
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
