package main

import "github.com/KTachibanaM/mear/agent"

func main() {
	err := agent.Agent()
	if err != nil {
		panic(err)
	}
}
