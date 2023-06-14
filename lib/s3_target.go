package lib

type S3Target struct {
	EndpointUrl     string
	Region          string
	BucketName      string
	ObjectKey       string
	AccessKeyId     string
	SecretAccessKey string
}
