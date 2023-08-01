package agent

import (
	"fmt"
	"path"

	"github.com/KTachibanaM/mear/internal/s3"
	"github.com/KTachibanaM/mear/internal/utils"
	log "github.com/sirupsen/logrus"
)

func Agent(agent_args *AgentArgs) error {
	log.Println("agent started")

	// 1. Download ffmpeg
	log.Println("downloading ffmpeg...")
	ffmpeg_workspace, err := GetWorkspaceDir("ffmpeg")
	if err != nil {
		return err
	}
	ffmpeg_executable, err := DownloadFfmpeg(ffmpeg_workspace)
	if err != nil {
		return err
	}

	// 2. Verify ffmpeg
	ffmpeg_version, err := GetFfmpegVersion(ffmpeg_executable)
	if err != nil {
		return err
	}
	log.Println(ffmpeg_version)

	// 2. Download videos
	video_workspace, err := GetWorkspaceDir("video")
	if err != nil {
		return err
	}
	var input_videos []string
	for i, job := range agent_args.Jobs {
		log.Printf("downloading video %v ...", job.S3Source)
		ext, err := utils.InferExt(job.S3Source.ObjectKey)
		if err != nil {
			return fmt.Errorf("could not infer the extension from the object key %v: %v", job.S3Source.ObjectKey, err)
		}
		input_video := path.Join(video_workspace, fmt.Sprintf("%v.in.%v", i, ext))
		err = s3.DownloadFile(input_video, job.S3Source, true)
		if err != nil {
			return err
		}
		input_videos = append(input_videos, input_video)
	}

	// 3. Convert video
	var output_videos []string
	for i, job := range agent_args.Jobs {
		log.Printf("converting video %v ...", job.S3Source)
		output_ext, err := utils.InferExt(job.S3Destination.ObjectKey)
		if err != nil {
			return err
		}
		output_video := path.Join(video_workspace, fmt.Sprintf("%v.out.%v", i, output_ext))
		err = ConvertVideo(ffmpeg_executable, input_videos[i], output_video, job.ExtraFfmpegArgs)
		if err != nil {
			return err
		}
		output_videos = append(output_videos, output_video)
	}

	// 4. Upload videos
	for i, job := range agent_args.Jobs {
		log.Printf("uploading video %v converted...", job.S3Source)
		err = s3.UploadFile(output_videos[i], job.S3Destination, true)
		if err != nil {
			return err
		}
	}

	return nil
}
