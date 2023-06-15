package internal

import (
	"bufio"
	"log"
	"os/exec"
)

func ConvertVideo(ffmpeg_executable, input_video, output_video string, extra_args []string) error {
	args := []string{"-i", input_video}
	args = append(args, extra_args...)
	args = append(args, output_video)

	cmd := exec.Command(ffmpeg_executable, args...)
	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		log.Println(m)
	}

	return cmd.Wait()
}
