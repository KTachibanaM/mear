package agent

import "github.com/KTachibanaM/mear/internal/s3"

type AgentJob struct {
	S3Source        *s3.S3Target `json:"s3Source"`
	S3Destination   *s3.S3Target `json:"s3Destination"`
	ExtraFfmpegArgs []string     `json:"extraFfmpegArgs"`
}

func NewAgentJob(s3_source, s3_destination *s3.S3Target, extraFfmpegargs []string) *AgentJob {
	return &AgentJob{
		S3Source:        s3_source,
		S3Destination:   s3_destination,
		ExtraFfmpegArgs: extraFfmpegargs,
	}
}

type AgentArgs struct {
	Jobs []*AgentJob `json:"jobs"`
}

func NewAgentArgs(s3_source, s3_destination *s3.S3Target, extraFfmpegargs []string) *AgentArgs {
	return &AgentArgs{
		Jobs: []*AgentJob{
			NewAgentJob(s3_source, s3_destination, extraFfmpegargs),
		},
	}
}

func NewAgentArgsWithoutJob() *AgentArgs {
	return &AgentArgs{
		Jobs: []*AgentJob{},
	}
}
