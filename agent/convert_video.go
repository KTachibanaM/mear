package agent

import "os/exec"

func ConvertVideo(ffmpeg_executable, input_video, output_video string) (string, error) {
	output, err := exec.Command(ffmpeg_executable, "-i", input_video, output_video).CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
