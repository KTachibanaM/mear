package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
)

func Agent(agent_args_s3_target *S3Target) error {
	log.Println("agent started")

	// 0. Download agent args
	log.Println("downloading agent args...")
	agent_args_bytes, err := ReadS3Target(agent_args_s3_target)
	if err != nil {
		return err
	}
	var agent_args AgentArgs
	err = json.Unmarshal(agent_args_bytes, &agent_args)
	if err != nil {
		return err
	}

	// 1. Download ffmpeg
	log.Println("downloading ffmpeg...")
	ffmpeg_workspace, err := os.MkdirTemp(os.TempDir(), "mear-ffmpeg-")
	if err != nil {
		return err
	}
	ffmpeg_executable, err := DownloadFfmpeg(ffmpeg_workspace)
	if err != nil {
		return err
	}
	log.Printf("ffmpeg is located at %s\n", ffmpeg_executable)

	// 2. Verify ffmpeg
	ffmpeg_version, err := GetFfmpegVersion(ffmpeg_executable)
	if err != nil {
		return err
	}
	log.Println(ffmpeg_version)

	// 2. Download video
	log.Println("downloading video...")
	video_workspace, err := os.MkdirTemp(os.TempDir(), "mear-video-")
	if err != nil {
		return err
	}
	input_video, err := DownloadVideo(video_workspace, agent_args.S3Source)
	if err != nil {
		return err
	}

	// 3. Convert video
	log.Println("converting video...")
	output_ext, err := InferExt(agent_args.S3Destination.ObjectKey)
	if err != nil {
		return err
	}
	output_video := path.Join(video_workspace, fmt.Sprintf("output.%s", output_ext))
	err = ConvertVideo(ffmpeg_executable, input_video, output_video, agent_args.ExtraFfmpegArgs)
	if err != nil {
		return err
	}

	// 4. Upload video
	log.Println("uploading video...")
	err = UploadVideo(output_video, agent_args.S3Destination)
	if err != nil {
		return err
	}

	log.Println("agent finished")
	return nil
}
