package agent

import "github.com/KTachibanaM/mear/lib"

type AgentArgs struct {
	S3Source        *lib.S3Target `json:"s3Source"`
	S3Destination   *lib.S3Target `json:"s3Destination"`
	ExtraFfmpegArgs []string      `json:"extraFfmpegArgs"`
}

func NewAgentArgs(s3_source, s3_destination *lib.S3Target, extraFfmpegArgs []string) *AgentArgs {
	return &AgentArgs{
		S3Source:        s3_source,
		S3Destination:   s3_destination,
		ExtraFfmpegArgs: extraFfmpegArgs,
	}
}
