package s3

type S3Session struct {
	EndpointUrl     string `json:"endpointUrl"`
	Region          string `json:"region"`
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	PathStyleUrl    bool   `json:"pathStyleUrl"`
}

func NewS3Session(endpoint_url, region, access_key_id, secret_access_key string, path_style_url bool) *S3Session {
	return &S3Session{
		EndpointUrl:     endpoint_url,
		Region:          region,
		AccessKeyId:     access_key_id,
		SecretAccessKey: secret_access_key,
		PathStyleUrl:    path_style_url,
	}
}
