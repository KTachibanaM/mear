package bucket

import (
	"fmt"

	"github.com/KTachibanaM/mear/internal/s3"
	"github.com/aws/aws-sdk-go/aws"

	aws_s3 "github.com/aws/aws-sdk-go/service/s3"
)

var source = s3.NewS3Target(
	"http://minio-source:9000",
	"us-east-1",
	"src",
	"MakeMine1948_256kb.rm",
	"minioadmin",
	"minioadmin",
	true,
)
var destination = s3.NewS3Target(
	"http://minio-destination:9000",
	"us-east-1",
	"dst",
	"output.mp4",
	"minioadmin",
	"minioadmin",
	true,
)
var logs = s3.NewS3Target(
	"http://minio-destination:9000",
	"us-east-1",
	"dst",
	"agent.log",
	"minioadmin",
	"minioadmin",
	true,
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

type DevcontainerBucketProvisioner struct{}

func NewDevcontainerBucketProvisioner() *DevcontainerBucketProvisioner {
	return &DevcontainerBucketProvisioner{}
}

func (p DevcontainerBucketProvisioner) Provision() (*s3.S3Target, *s3.S3Target, *s3.S3Target, error) {
	// Only delete destination objects before provisioning so that they can be retained for debugging after Teardown
	if err := deleteS3Target(destination); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to delete destination s3 target: %v", err)
	}
	if err := deleteS3Target(logs); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to delete logs s3 target: %v", err)
	}
	return source, destination, logs, nil
}

func (p DevcontainerBucketProvisioner) Teardown() error {
	// Actually do nothing for Teardown since destination objects are retained for debugging and will be deleted in the next provisioning
	return nil
}
