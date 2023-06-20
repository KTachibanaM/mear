package s3

import (
	"fmt"
	"io"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	mear_io "github.com/KTachibanaM/mear/internal/io"

	"github.com/KTachibanaM/mear/internal/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

// DownloadVideo downloads the video from S3 to the workspace_dir
// and returns the path to the downloaded file
func DownloadVideo(workspace_dir string, s3_target *S3Target) (string, error) {
	// Figure out the file extension
	ext, err := utils.InferExt(s3_target.ObjectKey)
	if err != nil {
		return "", fmt.Errorf("could not infer the extension from the object key %s: %w", s3_target.ObjectKey, err)
	}

	// Create the downloaded file
	downloaded := path.Join(workspace_dir, fmt.Sprintf("input.%s", ext))
	f, err := os.Create(downloaded)
	if err != nil {
		return "", fmt.Errorf("could not create the downloaded file: %w", err)
	}
	defer f.Close()

	// Create S3 session
	sess, err := CreateS3Session(s3_target)
	if err != nil {
		return "", fmt.Errorf("could not create S3 session for downloading video: %w", err)
	}

	// Check if video exists
	head_out, err := s3.New(sess).HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s3_target.BucketName),
		Key:    aws.String(s3_target.ObjectKey),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "NotFound" {
			return "", fmt.Errorf("video does not exist: %w", err)
		} else {
			return "", fmt.Errorf("could not check if video exists: %w", err)
		}
	}

	// Create request for downloading video
	req, err := s3.New(sess).GetObject(
		&s3.GetObjectInput{
			Bucket: aws.String(s3_target.BucketName),
			Key:    aws.String(s3_target.ObjectKey),
		},
	)
	if err != nil {
		return "", fmt.Errorf("could not create request for downloading video: %w", err)
	}

	// Download the video
	progress_writer := mear_io.NewProgressWriter(
		uint64(*head_out.ContentLength),
		func(progress float64) {
			log.Printf("downloaded %.2f%% of the video", progress)
		},
	)
	_, err = io.Copy(f, io.TeeReader(req.Body, progress_writer))
	if err != nil {
		return "", fmt.Errorf("could not download video: %w", err)
	}

	return downloaded, nil
}
