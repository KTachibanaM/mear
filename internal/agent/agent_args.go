package agent

import "github.com/KTachibanaM/mear/internal/s3"

type AgentJobDestination struct {
	S3Destination   *s3.S3Target `json:"s3Destination"`
	ExtraFfmpegArgs []string     `json:"extraFfmpegArgs"`
}

type AgentJob struct {
	S3Source        *s3.S3Target           `json:"s3Source"`
	JobDestinations []*AgentJobDestination `json:"jobDestinations"`
}

func NewAgentJob(s3_source *s3.S3Target) *AgentJob {
	return &AgentJob{
		S3Source:        s3_source,
		JobDestinations: []*AgentJobDestination{},
	}
}

func (job *AgentJob) AddJobDestination(s3_destination *s3.S3Target, extraFfmpegargs []string) {
	job.JobDestinations = append(job.JobDestinations, &AgentJobDestination{
		S3Destination:   s3_destination,
		ExtraFfmpegArgs: extraFfmpegargs,
	})
}

type AgentArgs struct {
	Jobs []*AgentJob `json:"jobs"`
}

func NewAgentArgs() *AgentArgs {
	return &AgentArgs{
		Jobs: []*AgentJob{},
	}
}

func (args *AgentArgs) AddJob(job *AgentJob) {
	args.Jobs = append(args.Jobs, job)
}
