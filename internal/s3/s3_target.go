package s3

import "fmt"

type S3Target struct {
	S3Bucket  *S3Bucket `json:"s3Bucket,omitempty"`
	ObjectKey string    `json:"objectKey"`
}

func (t *S3Target) String() string {
	return fmt.Sprintf("s3://%v/%v (endpoint=%v,region=%v)", t.S3Bucket.BucketName, t.ObjectKey, t.S3Bucket.S3Session.EndpointUrl, t.S3Bucket.S3Session.Region)
}

func NewS3Target(s3_bucket *S3Bucket, ObjectKey string) *S3Target {
	return &S3Target{
		S3Bucket:  s3_bucket,
		ObjectKey: ObjectKey,
	}
}
