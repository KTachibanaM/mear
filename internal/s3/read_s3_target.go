package s3

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

var Timeout = 5 * time.Second

// ReadS3Target reads an S3 target to a byte array.
func ReadS3Target(s3_target *S3Target) ([]byte, error) {
	// Create S3 session
	sess, err := CreateS3Session(s3_target)
	if err != nil {
		return nil, fmt.Errorf("could not create S3 session for reading: %w", err)
	}

	// Get object
	ctx, cancel := context.WithTimeout(context.TODO(), Timeout)
	defer cancel()

	obj, err := s3.New(sess).GetObjectWithContext(ctx, &s3.GetObjectInput{
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
