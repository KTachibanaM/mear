package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func CreateS3Session(s3_target *S3Target) (*session.Session, error) {
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
		return nil, err
	}
	return sess, nil
}
