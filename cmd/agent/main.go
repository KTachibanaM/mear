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
