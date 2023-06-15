package internal

import (
	"bufio"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func ConvertVideo(ffmpeg_executable, input_video, output_video string, extra_args []string) error {
	args := []string{"-i", input_video}
	args = append(args, extra_args...)
	args = append(args, output_video)

	cmd := exec.Command(ffmpeg_executable, args...)
	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	context_log := log.WithFields(log.Fields{
		"context": "ffmpeg",
	})

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		context_log.Println(m)
	}

	return cmd.Wait()
}
