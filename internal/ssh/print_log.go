package ssh

import (
	"encoding/json"

	"github.com/KTachibanaM/mear/internal/agent"
	log "github.com/sirupsen/logrus"
)

func print_agent_log_and_parse_failure(structured_log map[string]interface{}) *agent.AgentFailure {
	var parsed_agent_failure *agent.AgentFailure

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

		// Parse out potential failure
		failure, ok := structured_log["failure"].(map[string]interface{})
		if !ok {
			failure = nil
		} else {
			failure_bytes, _ := json.Marshal(failure)
			json.Unmarshal(failure_bytes, &parsed_agent_failure)
		}

		// Log agent log
		agent_log := log.WithFields(log.Fields{
			"context": context,
		})
		if failure != nil {
			agent_log = agent_log.WithFields(log.Fields{
				"failure": failure,
			})
		}
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

	return parsed_agent_failure
}

func print_log_and_parse_agent_failure(line string) *agent.AgentFailure {
	var structured_log map[string]interface{}
	err := json.Unmarshal([]byte(line), &structured_log)
	if err != nil {
		log.WithFields(log.Fields{
			"context": "ssh",
		}).Info(line)
		return nil
	}
	return print_agent_log_and_parse_failure(structured_log)
}
