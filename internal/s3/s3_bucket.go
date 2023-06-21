package s3

import "fmt"

type S3Bucket struct {
	S3Session  *S3Session `json:"s3Session"`
	BucketName string     `json:"bucketName"`
}

func (b *S3Bucket) String() string {
	return fmt.Sprintf("s3://%v (endpoint=%v,region=%v)", b.BucketName, b.S3Session.EndpointUrl, b.S3Session.Region)
}

func NewS3Bucket(s3_session *S3Session, bucket_name string) *S3Bucket {
	return &S3Bucket{
		S3Session:  s3_session,
		BucketName: bucket_name,
	}
}
