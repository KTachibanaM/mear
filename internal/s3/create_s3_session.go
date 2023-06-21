package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func CreateS3Session(s3_session *S3Session) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String(s3_session.EndpointUrl),
		Region:   aws.String(s3_session.Region),
		Credentials: credentials.NewStaticCredentials(
			s3_session.AccessKeyId,
			s3_session.SecretAccessKey,
			"",
		),
		S3ForcePathStyle: aws.Bool(s3_session.PathStyleUrl),
	})
	if err != nil {
		return nil, err
	}
	return sess, nil
}
