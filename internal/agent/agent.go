package agent

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

func Agent(agent_args *AgentArgs) *AgentFailure {
	log.Println("agent started")

	// 1. Download ffmpeg
	log.Println("downloading ffmpeg...")
	ffmpeg_workspace, err := GetWorkspaceDir("ffmpeg")
	if err != nil {
		return NewAgentFailureWithNonJobError(err)
	}
	ffmpeg_executable, err := DownloadFfmpeg(ffmpeg_workspace)
	if err != nil {
		return NewAgentFailureWithNonJobError(err)
	}

	// 2. Verify ffmpeg
	ffmpeg_version, err := GetFfmpegVersion(ffmpeg_executable)
	if err != nil {
		return NewAgentFailureWithNonJobError(err)
	}
	log.Println(ffmpeg_version)

	// 3. Execute agent jobs
	video_workspace, err := GetWorkspaceDir("video")
	if err != nil {
		return NewAgentFailureWithNonJobError(err)
	}
	failure := NewAgentFailure()
	for i, job := range agent_args.Jobs {
		log.Printf("executing agent job %v ...", i)
		err := ExecAgentJob(job, video_workspace, i, ffmpeg_executable)
		if err != nil {
			note := ""
			if strings.Contains(err.Error(), "signal: killed") {
				note = "ffmpeg might have been killed by os. you might want to use an engine with larger RAM."
			}
			failure.JobErrors = append(failure.JobErrors, NewAgentJobError(i, err, note))
		}
	}

	if len(failure.JobErrors) > 0 {
		return failure
	} else {
		return nil
	}
}
