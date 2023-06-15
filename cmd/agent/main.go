package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"

	"github.com/KTachibanaM/mear/internal"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: mear-agent <agent-args-s3-target-base64-encoded-json>")
	}

	// Base64 decode agent args S3 target JSON
	decoded, err := base64.StdEncoding.DecodeString(os.Args[1])
	if err != nil {
		log.Fatalln("failed to base64 decode agent args s3 target json: %w", err)
	}

	// Parse JSON
	var agent_args_s3_target internal.S3Target
	err = json.Unmarshal(decoded, &agent_args_s3_target)
	if err != nil {
		log.Fatalf("failed to parse agent args s3 target: %w", err)
	}

	// Run agent
	err = internal.Agent(&agent_args_s3_target)
	if err != nil {
		log.Fatalf("failed to run agent: %w", err)
	}
}
