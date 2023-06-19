package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/KTachibanaM/mear/internal"
	log "github.com/sirupsen/logrus"
)

var upload_log_interval = 10 * time.Second

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
	var agent_args internal.AgentArgs
	err = json.Unmarshal(decoded, &agent_args)
	if err != nil {
		log.Fatalf("failed to parse agent args: %v", err)
	}

	// Setup logger to log to both stdout and log file
	agent_workspace, err := internal.GetWorkspaceDir("agent")
	if err != nil {
		log.Fatalf("failed to create agent workspace: %v", err)
	}
	log_file := path.Join(agent_workspace, "agent.log")
	log_f, err := os.OpenFile(log_file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("failed to create log file: %v", err)
	}
	defer log_f.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, log_f))

	// Setup recurring job to upload log file to S3
	upload_ticker := time.NewTicker(upload_log_interval)
	go func() {
		for range upload_ticker.C {
			// TODO: should use fmt or log?
			fmt.Println("uploading log to s3...")
			err := internal.UploadFile(log_file, agent_args.S3Logs)
			if err != nil {
				// TODO: should use fmt or log?
				fmt.Printf("failed to upload log to s3: %v\n", err)
			}
		}
	}()

	// Run agent
	err = internal.Agent(&agent_args)
	if err != nil {
		log.WithFields(log.Fields{
			"result": false,
		}).Printf("failed to run agent: %v", err)
	} else {
		log.WithFields(log.Fields{
			"result": true,
		}).Info("successfully ran agent")
	}
	err = internal.UploadFile(log_file, agent_args.S3Logs)
	if err != nil {
		// TODO: should use fmt or log?
		fmt.Printf("failed to upload log to s3: %v\n", err)
	}
}
