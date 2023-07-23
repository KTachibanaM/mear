package ssh

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

func print_agent_log(structured_log map[string]interface{}) {
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
			agent_log.Error(msg)
		}
	}
}

func print_log(line string) {
	var structured_log map[string]interface{}
	err := json.Unmarshal([]byte(line), &structured_log)
	if err != nil {
		log.WithFields(log.Fields{
			"context": "ssh",
		}).Info(line)
	} else {
		print_agent_log(structured_log)
	}
}
