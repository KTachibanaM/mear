package agent

import (
	"fmt"
	"os"
	"path"
)

func Agent(agent_args *AgentArgs) error {
	// 1. Download ffmpeg
	ffmpeg_workspace, err := os.MkdirTemp(os.TempDir(), "mear-ffmpeg")
	if err != nil {
		return err
	}
	ffmpeg_executable, err := DownloadFfmpeg(ffmpeg_workspace)
	if err != nil {
		return err
	}
	println(ffmpeg_executable)

	// 2. Verify ffmpeg
	ffmpeg_version, err := RunFfmpegVersion(ffmpeg_executable)
	if err != nil {
		return err
	}
	println(ffmpeg_version)

	// 2. Download video
	video_workspace, err := os.MkdirTemp(os.TempDir(), "mear-video")
	if err != nil {
		return err
	}
	input_video, err := DownloadVideo(video_workspace, agent_args.S3Source)
	if err != nil {
		return err
	}

	// 3. Convert video
	output_ext, err := InferExt(agent_args.S3Destination.ObjectKey)
	if err != nil {
		return err
	}
	output_video := path.Join(video_workspace, fmt.Sprintf("output.%s", output_ext))
	ffmpeg_output, err := ConvertVideo(ffmpeg_executable, input_video, output_video, agent_args.ExtraFfmpegArgs)
	if err != nil {
		return err
	}
	println(ffmpeg_output)

	// 4. Upload video
	return UploadVideo(output_video, agent_args.S3Destination)
}
