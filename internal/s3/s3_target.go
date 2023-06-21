package s3

import "fmt"

type S3Target struct {
	EndpointUrl     string `json:"endpointUrl"`
	Region          string `json:"region"`
	BucketName      string `json:"bucketName"`
	ObjectKey       string `json:"objectName"`
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	PathStyleUrl    bool   `json:"pathStyleUrl"`
}

func (t *S3Target) String() string {
	return fmt.Sprintf("s3://%v/%v?endpoint=%v&region=%v", t.BucketName, t.ObjectKey, t.EndpointUrl, t.Region)
}

func NewS3Target(EndpointUrl, Region, BucketName, ObjectKey, AccessKeyId, SecretAccessKey string, PathStyleUrl bool) *S3Target {
	return &S3Target{
		EndpointUrl:     EndpointUrl,
		Region:          Region,
		BucketName:      BucketName,
		ObjectKey:       ObjectKey,
		AccessKeyId:     AccessKeyId,
		SecretAccessKey: SecretAccessKey,
		PathStyleUrl:    PathStyleUrl,
	}
}
