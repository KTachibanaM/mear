package agent

import (
	"fmt"
	"os"

	"github.com/KTachibanaM/mear/lib"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadVideo(video string, s3_target *lib.S3Target) error {
	// Open the video
	f, err := os.Open(video)
	if err != nil {
		return fmt.Errorf("could not open the video: %w", err)
	}

	// Create S3 session
	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String(s3_target.EndpointUrl),
		Region:   aws.String(s3_target.Region),
		Credentials: credentials.NewStaticCredentials(
			s3_target.AccessKeyId,
			s3_target.SecretAccessKey,
			"",
		),
		S3ForcePathStyle: aws.Bool(s3_target.PathStyleUrl),
	})
	if err != nil {
		return fmt.Errorf("could not create S3 session for download video: %w", err)
	}

	// Upload video
	_, err = s3manager.NewUploader(sess).Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3_target.BucketName),
		Key:    aws.String(s3_target.ObjectKey),
		Body:   f,
	})

	return err
}
