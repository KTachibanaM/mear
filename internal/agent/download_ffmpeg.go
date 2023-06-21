package agent

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	mear_io "github.com/KTachibanaM/mear/internal/io"
	log "github.com/sirupsen/logrus"
)

// DownloadFfmpeg downloads the latest version of ffmpeg for linux amd64
// It accepts a workspace_dir where the downloading and extracting is going to happen,
// and it returns the path to an ffmpeg executable
func DownloadFfmpeg(workspace_dir string) (string, error) {
	// https cert for ffmpeg website might be invalid?
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Skip verification of server certificate
			RootCAs:            nil,  // Use system default root CAs
		},
	}
	client := &http.Client{
		Transport: tr,
	}

	// Download the tar.xz file
	url := "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz"
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("could not download ffmpeg: %v", err)
	}
	defer resp.Body.Close()

	// Create the tar.gz file
	tar := path.Join(workspace_dir, "ffmpeg.tar.xz")
	out, err := os.Create(tar)
	if err != nil {
		return "", fmt.Errorf("could not create ffmpeg.tar.xz file: %v", err)
	}
	defer out.Close()

	// Get tar.gz file size
	var in io.Reader
	if resp.ContentLength != -1 {
		progress_writer := mear_io.NewProgressWriter(
			uint64(resp.ContentLength),
			func(progress float64) {
				log.Printf("downloaded %.2f%% of ffmpeg", progress)
			},
		)
		in = io.TeeReader(resp.Body, progress_writer)
	} else {
		log.Warn("could not get ffmpeg file size, progress will not be shown")
		in = resp.Body
	}

	// Write the body to file
	_, err = io.Copy(out, in)
	if err != nil {
		return "", fmt.Errorf("could not write ffmpeg.tar.xz file: %v", err)
	}

	// Untar the tar.gz file
	tar_f, err := os.Open(tar)
	if err != nil {
		return "", fmt.Errorf("could not open ffmpeg.tar.xz file: %v", err)
	}
	err = Untar(workspace_dir, tar_f)
	if err != nil {
		return "", fmt.Errorf("could not untar ffmpeg.tar.xz file: %v", err)
	}

	// Find the ffmpeg executable
	files, err := os.ReadDir(workspace_dir)
	if err != nil {
		return "", fmt.Errorf("could not read ffmpeg directory: %v", err)
	}
	var dirs []string
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		}
	}
	if len(dirs) != 1 {
		return "", fmt.Errorf("expecting one directory in ffmpeg directory, got %d", len(dirs))
	}
	ffmpeg_executable := path.Join(workspace_dir, dirs[0], "ffmpeg")

	return ffmpeg_executable, nil
}
