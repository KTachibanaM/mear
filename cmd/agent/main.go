package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/KTachibanaM/mear/internal"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mear-agent <path-to-json-agent-args>")
		os.Exit(1)
	}

	// Read JSON f
	f, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	// Parse JSON
	var args internal.AgentArgs
	err = json.Unmarshal(f, &args)
	if err != nil {
		panic(err)
	}

	// Run agent
	err = internal.Agent(&args)
	if err != nil {
		panic(err)
	}
}
