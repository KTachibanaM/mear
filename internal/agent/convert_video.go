package agent

import (
	"bufio"
	"bytes"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func ffmpegSplitLogLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\r'); i >= 0 {
		// Found a carriage return character, return the line up to that point
		return i + 1, data[0:i], nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// Found a newline character, return the line up to that point
		return i + 1, data[0:i], nil
	}
	if atEOF {
		// No more carriage return characters, return the remaining data as a line
		return len(data), data, nil
	}
	// Need more data
	return 0, nil, nil
}

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
	scanner.Split(ffmpegSplitLogLines)
	for scanner.Scan() {
		m := scanner.Text()
		context_log.Println(m)
	}

	return cmd.Wait()
}
