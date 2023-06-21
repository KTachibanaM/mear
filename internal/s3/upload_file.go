package s3

import (
	"fmt"
	"io"
	"os"

	mear_io "github.com/KTachibanaM/mear/internal/io"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
)

func UploadFile(file string, s3_target *S3Target, print_progress bool) error {
	// Open the file
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("could not open the file: %v", err)
	}
	defer f.Close()

	// Get the file stat
	f_stat, err := f.Stat()
	if err != nil {
		return fmt.Errorf("could not get the file stat: %v", err)
	}

	// Create S3 session
	sess, err := CreateS3Session(s3_target.S3Bucket.S3Session)
	if err != nil {
		return fmt.Errorf("could not create S3 session for uploading file: %v", err)
	}

	// Upload file
	progress_writer := mear_io.NewProgressWriter(
		uint64(f_stat.Size()),
		func(progress float64) {
			if print_progress {
				log.Printf("uploaded %.2f%% of the file", progress)
			}
		},
	)
	_, err = s3manager.NewUploader(sess).Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3_target.S3Bucket.BucketName),
		Key:    aws.String(s3_target.ObjectKey),
		Body:   io.TeeReader(f, progress_writer),
	})
	if err != nil {
		return fmt.Errorf("could not upload file: %v", err)
	}

	return nil
}
