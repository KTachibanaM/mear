package agent

import "fmt"

type AgentJobError struct {
	JobIndex int    `json:"jobIndex"`
	Error    string `json:"error"`
	Note     string `json:"note"`
}

func NewAgentJobError(jobIndex int, error error, note string) *AgentJobError {
	return &AgentJobError{
		JobIndex: jobIndex,
		Error:    error.Error(),
		Note:     note,
	}
}

type AgentFailure struct {
	NonJobError string           `json:"nonJobError"`
	JobErrors   []*AgentJobError `json:"jobErrors"`
}

func (a *AgentFailure) FirstJobError() error {
	if a.NonJobError != "" {
		return fmt.Errorf(a.NonJobError)
	}
	if len(a.JobErrors) > 0 {
		return fmt.Errorf(a.JobErrors[0].Error)
	}
	return nil
}

func NewAgentFailure() *AgentFailure {
	return &AgentFailure{
		NonJobError: "",
		JobErrors:   make([]*AgentJobError, 0),
	}
}

func NewAgentFailureWithNonJobError(err error) *AgentFailure {
	return &AgentFailure{
		NonJobError: err.Error(),
		JobErrors:   make([]*AgentJobError, 0),
	}
}
