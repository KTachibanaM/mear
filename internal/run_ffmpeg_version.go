package internal

import "os/exec"

// RunFfmpegVersion runs ffmpeg -version and returns the output
func RunFfmpegVersion(ffmpeg_executable_path string) (string, error) {
	output, err := exec.Command(ffmpeg_executable_path, "-version").CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
