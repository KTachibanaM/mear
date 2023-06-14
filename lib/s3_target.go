package lib

type S3Target struct {
	EndpointUrl     string
	Region          string
	BucketName      string
	ObjectKey       string
	AccessKeyId     string
	SecretAccessKey string
	PathStyleUrl    bool
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
