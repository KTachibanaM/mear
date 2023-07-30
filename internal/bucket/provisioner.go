package bucket

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	mear_s3 "github.com/KTachibanaM/mear/internal/s3"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-sdk-go/aws"
)

func ProvisionBucket(s3_bucket *mear_s3.S3Bucket) error {
	sess, err := mear_s3.CreateS3Session(s3_bucket.S3Session)
	if err != nil {
		return fmt.Errorf("could not create S3 session for bucket %v provisioning: %v", s3_bucket, err)
	}

	log.Printf("creating bucket %v ...", s3_bucket)
	_, err = s3.New(sess).CreateBucket(
		&s3.CreateBucketInput{
			Bucket: aws.String(s3_bucket.BucketName),
		},
	)
	if err != nil {
		return fmt.Errorf("could not create bucket %v: %v", s3_bucket, err)
	}
	return nil
}

func TeardownBucket(s3_bucket *mear_s3.S3Bucket) error {
	sess, err := mear_s3.CreateS3Session(s3_bucket.S3Session)
	if err != nil {
		return fmt.Errorf("could not create S3 session for bucket %v teardown: %v", s3_bucket, err)
	}

	log.Printf("emptying bucket %v ...", s3_bucket)
	resp, err := s3.New(sess).ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(s3_bucket.BucketName),
	})
	if err != nil {
		return fmt.Errorf("could not list objects in bucket %v for deletion: %v", s3_bucket, err)
	}
	var keys []*s3.ObjectIdentifier
	for _, obj := range resp.Contents {
		keys = append(keys, &s3.ObjectIdentifier{
			Key: obj.Key,
		})
	}
	if len(keys) > 0 {
		_, err = s3.New(sess).DeleteObjects(&s3.DeleteObjectsInput{
			Bucket: aws.String(s3_bucket.BucketName),
			Delete: &s3.Delete{
				Objects: keys,
				Quiet:   aws.Bool(true),
			},
		})
		if err != nil {
			return fmt.Errorf("could not delete objects in bucket %v: %v", s3_bucket, err)
		}
	}

	log.Printf("deleting bucket %v ...", s3_bucket)
	_, err = s3.New(sess).DeleteBucket(
		&s3.DeleteBucketInput{
			Bucket: aws.String(s3_bucket.BucketName),
		},
	)
	if err != nil {
		return fmt.Errorf("could not delete bucket %v: %v", s3_bucket, err)
	}

	return nil
}
