package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/KTachibanaM/mear/agent"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mear <path-to-json-agent-args>")
		os.Exit(1)
	}

	// Read JSON f
	f, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	// Parse JSON
	var args agent.AgentArgs
	err = json.Unmarshal(f, &args)
	if err != nil {
		panic(err)
	}

	// Run agent
	err = agent.Agent(&args)
	if err != nil {
		panic(err)
	}
}
