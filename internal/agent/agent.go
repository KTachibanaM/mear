package agent

import (
	"fmt"
	"path"
	"strings"

	"github.com/KTachibanaM/mear/internal/s3"
	"github.com/KTachibanaM/mear/internal/utils"
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
	for jobIndex, job := range agent_args.Jobs {
		// 1. Download video
		ext, err := utils.InferExt(job.S3Source.ObjectKey)
		if err != nil {
			failure.JobErrors = append(failure.JobErrors, NewAgentJobError(jobIndex, 0, err, ""))
			break
		}
		log.Printf("downloading video %v ...", job.S3Source)
		input_video := path.Join(video_workspace, fmt.Sprintf("%v.in.%v", jobIndex, ext))
		err = s3.DownloadFile(input_video, job.S3Source, true)
		if err != nil {
			failure.JobErrors = append(failure.JobErrors, NewAgentJobError(jobIndex, 0, err, ""))
			break
		}

		for jobDestinationIndex, jobDestination := range job.JobDestinations {
			log.Printf("executing agent job %v for destination %v ...", jobIndex, jobDestinationIndex)
			err := ExecAgentJobForDestination(input_video, jobDestination, video_workspace, jobIndex, jobDestinationIndex, ffmpeg_executable)
			if err != nil {
				note := ""
				if strings.Contains(err.Error(), "signal: killed") {
					note = "ffmpeg might have been killed by os. you might want to use an engine with larger RAM."
				}
				failure.JobErrors = append(failure.JobErrors, NewAgentJobError(jobIndex, jobDestinationIndex, err, note))
			}
		}
	}

	if len(failure.JobErrors) > 0 {
		return failure
	} else {
		return nil
	}
}
