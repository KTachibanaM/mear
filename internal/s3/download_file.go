package s3

import (
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"

	mear_io "github.com/KTachibanaM/mear/internal/io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

func DownloadFile(file string, s3_target *S3Target, print_progress bool) error {
	// Create the downloaded file
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("could not create the downloaded file: %v", err)
	}
	defer f.Close()

	// Create S3 session
	sess, err := CreateS3Session(s3_target.S3Bucket.S3Session)
	if err != nil {
		return fmt.Errorf("could not create S3 session for downloading file: %v", err)
	}

	// Check if the file exists
	head_out, err := s3.New(sess).HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s3_target.S3Bucket.BucketName),
		Key:    aws.String(s3_target.ObjectKey),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "NotFound" {
			return fmt.Errorf("file does not exist: %v", err)
		} else {
			return fmt.Errorf("could not check if file exists: %v", err)
		}
	}

	// Create request for downloading file
	req, err := s3.New(sess).GetObject(
		&s3.GetObjectInput{
			Bucket: aws.String(s3_target.S3Bucket.BucketName),
			Key:    aws.String(s3_target.ObjectKey),
		},
	)
	if err != nil {
		return fmt.Errorf("could not create request for downloading file: %v", err)
	}

	// Download the file
	progress_writer := mear_io.NewProgressWriter(
		uint64(*head_out.ContentLength),
		func(progress float64) {
			if print_progress {
				log.Printf("downloaded %.2f%% of the file", progress)
			}
		},
	)
	_, err = io.Copy(f, io.TeeReader(req.Body, progress_writer))
	if err != nil {
		return fmt.Errorf("could not download file: %v", err)
	}

	return nil
}
