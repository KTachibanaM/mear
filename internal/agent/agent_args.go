package agent

import "github.com/KTachibanaM/mear/internal/s3"

type AgentArgs struct {
	S3Source        *s3.S3Target `json:"s3Source"`
	S3Destination   *s3.S3Target `json:"s3Destination"`
	S3Logs          *s3.S3Target `json:"s3Logs"`
	ExtraFfmpegArgs []string     `json:"extraFfmpegArgs"`
}
