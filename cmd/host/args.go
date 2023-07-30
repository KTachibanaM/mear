package main

import "github.com/KTachibanaM/mear/internal/host"

type HostArgs struct {
	Jobs                         []*host.Job `json:"jobs"`
	AgentExecutionTimeoutMinutes int         `json:"agent_execution_timeout_minutes"`
	Stack                        string      `json:"stack"`
	DropletRam                   int         `json:"droplet_ram"`
	DropletCpu                   int         `json:"droplet_cpu"`
	DoAccessKeyId                string      `json:"do_access_key_id"`
	DoSecretAccessKey            string      `json:"do_secret_access_key"`
	DoToken                      string      `json:"do_token"`
}

type HostResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func NewHostResult(success bool, message string) *HostResult {
	return &HostResult{
		Success: success,
		Message: message,
	}
}
