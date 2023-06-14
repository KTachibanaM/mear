package agent

import "os/exec"

func ConvertVideo(ffmpeg_executable, input_video, output_video string, extra_args []string) (string, error) {
	var args = []string{"-i", input_video}
	for _, arg := range extra_args {
		args = append(args, arg)
	}
	args = append(args, output_video)

	output, err := exec.Command(ffmpeg_executable, args...).CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
