package internal

// ProgressWriter counts the number of bytes written to it
// and reports progress on every 10% bytes written using the given callback
// It implements to the io.Writer interface and we can pass this into io.TeeReader()
type ProgressWriter struct {
	Expected          uint64
	LogCallback       func(float64)
	counted           uint64
	previous_progress float64
}

func (wc *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	wc.counted += uint64(n)
	cur_progress := float64(wc.counted) / float64(wc.Expected) * 100
	if cur_progress-wc.previous_progress >= 10 {
		wc.LogCallback(cur_progress)
		wc.previous_progress = cur_progress
	}
	return n, nil
}

func NewProgressWriter(expected_bytes uint64, log_callback func(float64)) *ProgressWriter {
	return &ProgressWriter{
		Expected:    expected_bytes,
		LogCallback: log_callback,
	}
}
