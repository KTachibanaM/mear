package bucket

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	mear_s3 "github.com/KTachibanaM/mear/internal/s3"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-sdk-go/aws"
)

type BucketProvisioner struct {
	s3_session  *mear_s3.S3Session
	bucket_name string
}

func NewBucketProvisioner(
	s3_session *mear_s3.S3Session,
) *BucketProvisioner {
	return &BucketProvisioner{
		s3_session:  s3_session,
		bucket_name: "",
	}
}

func (p *BucketProvisioner) Provision(bucket_name string) (*mear_s3.S3Bucket, error) {
	bucket := mear_s3.NewS3Bucket(p.s3_session, bucket_name)

	sess, err := mear_s3.CreateS3Session(p.s3_session)
	if err != nil {
		return nil, fmt.Errorf("could not create S3 session for bucket %v provisioning: %v", bucket, err)
	}

	log.Printf("creating bucket %v...", bucket)
	_, err = s3.New(sess).CreateBucket(
		&s3.CreateBucketInput{
			Bucket: aws.String(bucket_name),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("could not create bucket %v: %v", bucket, err)
	}
	p.bucket_name = bucket_name
	return bucket, nil

}

func (p *BucketProvisioner) Teardown() error {
	if p.bucket_name == "" {
		return fmt.Errorf("bucket '%v' was never provisioned", p.bucket_name)
	}

	bucket := mear_s3.NewS3Bucket(p.s3_session, p.bucket_name)

	sess, err := mear_s3.CreateS3Session(p.s3_session)
	if err != nil {
		return fmt.Errorf("could not create S3 session for bucket %v teardown: %v", bucket, err)
	}

	log.Printf("emptying bucket %v...", bucket)
	resp, err := s3.New(sess).ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(p.bucket_name),
	})
	if err != nil {
		return fmt.Errorf("could not list objects in bucket %v for deletion: %v", bucket, err)
	}
	var keys []*s3.ObjectIdentifier
	for _, obj := range resp.Contents {
		keys = append(keys, &s3.ObjectIdentifier{
			Key: obj.Key,
		})
	}
	if len(keys) > 0 {
		_, err = s3.New(sess).DeleteObjects(&s3.DeleteObjectsInput{
			Bucket: aws.String(p.bucket_name),
			Delete: &s3.Delete{
				Objects: keys,
				Quiet:   aws.Bool(true),
			},
		})
		if err != nil {
			return fmt.Errorf("could not delete objects in bucket %v: %v", bucket, err)
		}
	}

	log.Printf("deleting bucket %v...", bucket)
	_, err = s3.New(sess).DeleteBucket(
		&s3.DeleteBucketInput{
			Bucket: aws.String(p.bucket_name),
		},
	)
	if err != nil {
		return fmt.Errorf("could not delete bucket %v: %v", bucket, err)
	}

	return nil
}
