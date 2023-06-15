package internal

import "os/exec"

func ConvertVideo(ffmpeg_executable, input_video, output_video string, extra_args []string) (string, error) {
	args := []string{"-i", input_video}
	args = append(args, extra_args...)
	args = append(args, output_video)

	output, err := exec.Command(ffmpeg_executable, args...).CombinedOutput()
	if err != nil {
		println("failed to run ffmpeg, &w", err)
		return "", err
	}
	return string(output), nil
}
