package internal

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadVideo(video string, s3_target *S3Target) error {
	// Open the video
	f, err := os.Open(video)
	if err != nil {
		return fmt.Errorf("could not open the video: %w", err)
	}

	// Create S3 session
	sess, err := CreateS3Session(s3_target)
	if err != nil {
		return fmt.Errorf("could not create S3 session for uploading video: %w", err)
	}

	// Upload video
	_, err = s3manager.NewUploader(sess).Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3_target.BucketName),
		Key:    aws.String(s3_target.ObjectKey),
		Body:   f,
	})

	return err
}
