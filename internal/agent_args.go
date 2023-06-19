package internal

type AgentArgs struct {
	S3Source        *S3Target `json:"s3Source"`
	S3Destination   *S3Target `json:"s3Destination"`
	S3Logs          *S3Target `json:"s3Logs"`
	ExtraFfmpegArgs []string  `json:"extraFfmpegArgs"`
}
