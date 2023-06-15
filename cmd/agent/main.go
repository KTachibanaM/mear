package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/KTachibanaM/mear/internal"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mear-agent <agent-args-s3-target-base64-encoded-json>")
		os.Exit(1)
	}

	// Base64 decode agent args S3 target JSON
	decoded, err := base64.StdEncoding.DecodeString(os.Args[1])
	if err != nil {
		panic(err)
	}

	// Parse JSON
	var agent_args_s3_target internal.S3Target
	err = json.Unmarshal(decoded, &agent_args_s3_target)
	if err != nil {
		panic(err)
	}

	// Run agent
	err = internal.Agent(&agent_args_s3_target)
	if err != nil {
		panic(err)
	}
}
