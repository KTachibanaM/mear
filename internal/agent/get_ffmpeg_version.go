package agent

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetFfmpegVersion runs ffmpeg -version and returns the first line of its output
func GetFfmpegVersion(ffmpeg_executable_path string) (string, error) {
	output, err := exec.Command(ffmpeg_executable_path, "-version").CombinedOutput()
	if err != nil {
		return "", err
	}
	splits := strings.Split(string(output), "\n")
	if len(splits) == 0 {
		return "", fmt.Errorf("could not get any line of ffmpeg -version output")
	}
	return splits[0], nil
}
