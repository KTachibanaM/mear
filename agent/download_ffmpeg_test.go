package agent

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadFFmpeg(t *testing.T) {
	testing_ffmpeg_workspace, err := os.MkdirTemp(os.TempDir(), "testing-mear-ffmpeg")
	assert.NoError(t, err)
	ffmpeg_executable, err := DownloadFfmpeg(testing_ffmpeg_workspace)
	assert.NoError(t, err)
	version_string, err := RunFfmpegVersion(ffmpeg_executable)
	assert.NoError(t, err)
	assert.Contains(t, version_string, "ffmpeg version")
}
