package engine

type EngineProvisioner interface {
	Provision(agent_binary_url, encoded_agent_args string) error
	Teardown() error
}
