package agent

import (
	"fmt"
	"os"
	"path"

	"github.com/KTachibanaM/mear/lib"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// DownloadVideo downloads the video from S3 to the workspace_dir
// and returns the path to the downloaded file
func DownloadVideo(workspace_dir string, s3_target *lib.S3Target) (string, error) {
	// Figure out the file extension
	ext, err := InferExt(s3_target.ObjectKey)
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
		return "", fmt.Errorf("could not create S3 session for download video: %w", err)
	}

	// Check if video exists
	_, err = s3.New(sess).HeadObject(&s3.HeadObjectInput{
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

	// Download video
	_, err = s3manager.NewDownloader(sess).Download(
		f, &s3.GetObjectInput{
			Bucket: aws.String(s3_target.BucketName),
			Key:    aws.String(s3_target.ObjectKey),
		},
	)
	if err != nil {
		return "", fmt.Errorf("could not download video: %w", err)
	}

	return downloaded, nil
}
