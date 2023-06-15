package internal

type AgentArgs struct {
	S3Source        *S3Target `json:"s3Source"`
	S3Destination   *S3Target `json:"s3Destination"`
	ExtraFfmpegArgs []string  `json:"extraFfmpegArgs"`
}

func NewAgentArgs(s3_source, s3_destination *S3Target, extraFfmpegArgs []string) *AgentArgs {
	return &AgentArgs{
		S3Source:        s3_source,
		S3Destination:   s3_destination,
		ExtraFfmpegArgs: extraFfmpegArgs,
	}
}
