package agent

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
