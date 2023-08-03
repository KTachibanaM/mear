package agent

import (
	"fmt"
	"path"

	"github.com/KTachibanaM/mear/internal/s3"
	"github.com/KTachibanaM/mear/internal/utils"
	log "github.com/sirupsen/logrus"
)

func ExecAgentJob(job *AgentJob, video_workspace string, id int, ffmpeg_executable string) error {
	// 1. Download video
	log.Printf("downloading video %v ...", job.S3Source)
	ext, err := utils.InferExt(job.S3Source.ObjectKey)
	if err != nil {
		return fmt.Errorf("could not infer the extension from the object key %v: %v", job.S3Source.ObjectKey, err)
	}
	input_video := path.Join(video_workspace, fmt.Sprintf("%v.in.%v", id, ext))
	err = s3.DownloadFile(input_video, job.S3Source, true)
	if err != nil {
		return err
	}

	// 2. Convert video
	log.Printf("converting video %v ...", job.S3Source)
	output_ext, err := utils.InferExt(job.S3Destination.ObjectKey)
	if err != nil {
		return err
	}
	output_video := path.Join(video_workspace, fmt.Sprintf("%v.out.%v", id, output_ext))
	err = ConvertVideo(ffmpeg_executable, input_video, output_video, job.ExtraFfmpegArgs)
	if err != nil {
		return err
	}

	// 3. Upload video
	log.Printf("uploading converted video %v ...", job.S3Source)
	err = s3.UploadFile(output_video, job.S3Destination, true)
	if err != nil {
		return err
	}

	return nil
}
