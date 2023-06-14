package lib

type S3Target struct {
	EndpointUrl     string `json:"endpointUrl"`
	Region          string `json:"region"`
	BucketName      string `json:"bucketName"`
	ObjectKey       string `json:"objectName"`
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	PathStyleUrl    bool   `json:"pathStyleUrl"`
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
