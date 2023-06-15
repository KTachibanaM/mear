package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"path"

	"github.com/KTachibanaM/mear/internal"
	log "github.com/sirupsen/logrus"
)

func main() {
	// setup logger to log to both stdout and log file
	log.SetFormatter(&log.JSONFormatter{})
	agent_workspace, err := os.MkdirTemp(os.TempDir(), "mear-agent-")
	if err != nil {
		log.Fatalf("failed to create agent workspace: %v", err)
	}
	log_file := path.Join(agent_workspace, "mear-agent.log")
	log_f, err := os.OpenFile(log_file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("failed to create log file: %v", err)
	}
	defer log_f.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, log_f))

	if len(os.Args) < 2 {
		log.Fatalln("Usage: mear-agent <agent-args-json-base64-encoded>")
	}

	// Base64 decode agent args
	decoded, err := base64.StdEncoding.DecodeString(os.Args[1])
	if err != nil {
		log.Fatalf("failed to base64 decode agent args: %v", err)
	}

	// Parse JSON
	var agent_args internal.AgentArgs
	err = json.Unmarshal(decoded, &agent_args)
	if err != nil {
		log.Fatalf("failed to parse agent args: %v", err)
	}

	// Run agent
	err = internal.Agent(&agent_args)
	if err != nil {
		log.Fatalf("failed to run agent: %v", err)
	}
}
