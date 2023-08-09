package agent

import (
	"fmt"
	"path"

	"github.com/KTachibanaM/mear/internal/s3"
	"github.com/KTachibanaM/mear/internal/utils"
	log "github.com/sirupsen/logrus"
)

func ExecAgentJobForDestination(input_video string, job *AgentJobDestination, video_workspace string, jobIndex, jobDestinationIndex int, ffmpeg_executable string) error {
	// 2. Convert video
	log.Printf("converting video %v ...", input_video)
	output_ext, err := utils.InferExt(job.S3Destination.ObjectKey)
	if err != nil {
		return err
	}
	output_video := path.Join(video_workspace, fmt.Sprintf("%v.%v.out.%v", jobIndex, jobDestinationIndex, output_ext))
	err = ConvertVideo(ffmpeg_executable, input_video, output_video, job.ExtraFfmpegArgs)
	if err != nil {
		return err
	}

	// 3. Upload video
	log.Printf("uploading converted video %v ...", output_video)
	err = s3.UploadFile(output_video, job.S3Destination, true)
	if err != nil {
		return err
	}

	return nil
}
