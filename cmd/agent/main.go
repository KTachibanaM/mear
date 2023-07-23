package main

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/KTachibanaM/mear/internal/agent"
	log "github.com/sirupsen/logrus"
)

var UploadLogInterval = 10 * time.Second

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	if len(os.Args) < 2 {
		log.Fatalln("Usage: mear-agent <agent-args-json-base64-encoded>")
	}

	// Base64 decode agent args
	decoded, err := base64.StdEncoding.DecodeString(os.Args[1])
	if err != nil {
		log.Fatalf("failed to base64 decode agent args: %v", err)
	}

	// Parse JSON
	var agent_args agent.AgentArgs
	err = json.Unmarshal(decoded, &agent_args)
	if err != nil {
		log.Fatalf("failed to parse agent args: %v", err)
	}

	// Run agent
	err = agent.Agent(&agent_args)
	if err != nil {
		if strings.Contains(err.Error(), "signal: killed") {
			log.Warnln("ffmpeg might have been killed by os. you might want to use an engine with larger RAM.")
		}
		log.Fatalf("failed to run agent: %v", err)
	} else {
		log.Info("successfully ran agent")
	}
}
