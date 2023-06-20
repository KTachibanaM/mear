package s3

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ReadS3Target reads an S3 target to a byte array.
func ReadS3Target(s3_target *S3Target) ([]byte, error) {
	// Create S3 session
	sess, err := CreateS3Session(s3_target)
	if err != nil {
		return nil, fmt.Errorf("could not create S3 session for reading: %w", err)
	}

	// Get object
	obj, err := s3.New(sess).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3_target.BucketName),
		Key:    aws.String(s3_target.ObjectKey),
	})
	if err != nil {
		return nil, fmt.Errorf("could not get object for reading: %w", err)
	}

	// Read object contents
	contents, err := io.ReadAll(obj.Body)
	if err != nil {
		return nil, err
	}
	defer obj.Body.Close()

	return contents, nil
}
