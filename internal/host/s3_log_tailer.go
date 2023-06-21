package host

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/KTachibanaM/mear/internal/s3"
	log "github.com/sirupsen/logrus"
)

// Tailing for 2 hours max
var TailingMaxIntervals = 720
var TailingInterval = 10 * time.Second

var HeartbeatWarnTolerance = 30 * time.Second
var HeartbeatErrorTolerance = 5 * time.Minute

type S3LogsTailer struct {
	next_log_position int
	s3_logs           *s3.S3Target
	last_heartbeat    time.Time
}

func NewS3LogsTailer(s3_logs *s3.S3Target) *S3LogsTailer {
	return &S3LogsTailer{
		next_log_position: 0,
		s3_logs:           s3_logs,
		last_heartbeat:    time.Now(),
	}
}

// processStructuredLog processes structured log
// It returns a pair of booleans
// The first is true if the agent run is finished
// The second is true if the agent run is successful
func (t *S3LogsTailer) processStructuredLog(structured_log map[string]interface{}) (bool, bool) {
	// Parse out msg
	msg, ok := structured_log["msg"].(string)
	if !ok {
		log.Warnf("failed to find msg in structured log: %v", structured_log)
	} else {
		// Parse out level
		level, ok := structured_log["level"].(string)
		if !ok {
			log.Warnf("failed to find level in structured log: %v", structured_log)
			level = "info"
		}

		// Parse out potential context
		context, ok := structured_log["context"].(string)
		if !ok {
			context = "agent"
		}

		// Parse out potential heartbeat
		heartbeat, ok := structured_log["heartbeat"].(bool)
		if !ok {
			heartbeat = false
		}

		if heartbeat {
			// Log heartbeat
			log.Infoln("received heartbeat from agent")
			t.last_heartbeat = time.Now()
		} else {
			// Log agent log
			agent_log := log.WithFields(log.Fields{
				"context": context,
			})
			switch level {
			case "trace":
				agent_log.Trace(msg)
			case "debug":
				agent_log.Debug(msg)
			case "info":
				agent_log.Info(msg)
			case "warn":
				agent_log.Warn(msg)
			case "error":
				agent_log.Error(msg)
			case "fatal":
				agent_log.Fatal(msg)
			}
		}
	}

	// Parse out result
	result, ok := structured_log["result"].(bool)
	if ok {
		return true, result
	} else {
		return false, false
	}
}

// processLogLines parses log lines as structured logs and processes them
// It returns a pair of booleans
// The first is true if the agent run is finished
// The second is true if the agent run is successful
func (t *S3LogsTailer) processLogLines(log_lines []string) (bool, bool) {
	for _, log_line := range log_lines {
		var structured_log map[string]interface{}
		err := json.Unmarshal([]byte(log_line), &structured_log)
		if err != nil {
			log.Warnf("failed to unmarshal log line: %v", log_line)
		} else {
			terminate, result := t.processStructuredLog(structured_log)
			if terminate {
				return true, result
			}
		}
	}
	return false, false
}

// interval reads logs from s3_logs and processes them
// It returns a pair of booleans
// The first is true if the agent run is finished
// The second is true if the agent run is successful
func (t *S3LogsTailer) interval() (bool, bool) {
	// Process heartbeat
	now := time.Now()
	if now.Sub(t.last_heartbeat) > HeartbeatErrorTolerance {
		log.Errorf("last heartbeat was %v seconds ago, terminating...", now.Sub(t.last_heartbeat).Seconds())
		return true, false
	} else if now.Sub(t.last_heartbeat) > HeartbeatWarnTolerance {
		log.Warnf("last heartbeat was late (%v seconds ago)", now.Sub(t.last_heartbeat).Seconds())
	}

	// Process logs
	bytes, err := s3.ReadS3Target(t.s3_logs)
	if err != nil {
		log.Warnf("failed to read logs, wait for next interval: %v", err)
		return false, false
	}
	logs := string(bytes)
	raw_log_lines := strings.Split(logs, "\n")
	log_lines := []string{}
	for _, raw_log_line := range raw_log_lines {
		trimmed_log_line := strings.Trim(raw_log_line, " ")
		if trimmed_log_line != "" {
			log_lines = append(log_lines, trimmed_log_line)
		}
	}
	if t.next_log_position >= len(log_lines) {
		return false, false
	}
	terminate, result := t.processLogLines(log_lines[t.next_log_position:])
	if terminate {
		return true, result
	}
	t.next_log_position = len(log_lines) + 1
	return false, false
}

// Tail tails logs from s3_logs and returns whether the agent run is successful
func (t *S3LogsTailer) Tail() bool {
	for i := 0; i < TailingMaxIntervals; i++ {
		terminate, result := t.interval()
		if terminate {
			return result
		}
		time.Sleep(TailingInterval)
	}

	return false
}
