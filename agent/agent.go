package agent

import (
	"os"
)

// func Agent(input, output *lib.S3Target) error {
func Agent() error {
	// 1. Download ffmpeg
	ffmpeg_workspace, err := os.MkdirTemp(os.TempDir(), "mear-ffmpeg")
	if err != nil {
		return err
	}
	ffmpeg_executable, err := DownloadFfmpeg(ffmpeg_workspace)
	if err != nil {
		return err
	}

	// 2. Verify ffmpeg
	ffmpeg_version, err := RunFfmpegVersion(ffmpeg_executable)
	if err != nil {
		return err
	}
	println(ffmpeg_version)

	// 2. Download video
	// video_workspace, err := os.MkdirTemp(os.TempDir(), "mear-video")
	// if err != nil {
	// 	return err
	// }
	// input_video, err := DownloadVideo(video_workspace, input)
	// if err != nil {
	// 	return err
	// }

	// 3. Convert video

	// 4. Upload video
	return nil
}
