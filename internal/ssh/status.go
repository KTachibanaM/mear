package ssh

import "github.com/KTachibanaM/mear/internal/agent"

const (
	SshSuccess = iota
	SshError
	SshAgentFailure
)

type SshStatus struct {
	Status       int
	Err          error
	AgentFailure *agent.AgentFailure
}

func NewSshSuccess() *SshStatus {
	return &SshStatus{
		Status: SshSuccess,
	}
}

func NewSshError(err error) *SshStatus {
	return &SshStatus{
		Status: SshError,
		Err:    err,
	}
}

func NewSshAgentFailure(agent_failure *agent.AgentFailure) *SshStatus {
	return &SshStatus{
		Status:       SshAgentFailure,
		AgentFailure: agent_failure,
	}
}
