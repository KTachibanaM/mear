package agent

import "github.com/KTachibanaM/mear/internal/s3"

type AgentArgs struct {
	S3Source        *s3.S3Target `json:"s3Source"`
	S3Destination   *s3.S3Target `json:"s3Destination"`
	S3Logs          *s3.S3Target `json:"s3Logs"`
	ExtraFfmpegArgs []string     `json:"extraFfmpegArgs"`
}

func NewAgentArgs(s3_source, s3_destination, s3_logs *s3.S3Target, extraFfmpegargs []string) *AgentArgs {
	return &AgentArgs{
		S3Source:        s3_source,
		S3Destination:   s3_destination,
		S3Logs:          s3_logs,
		ExtraFfmpegArgs: extraFfmpegargs,
	}
}