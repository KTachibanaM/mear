package engine

type EngineProvisioner interface {
	Provision(agent_binary_url string, ssh_public_key []byte) (string, error)
	Teardown() error
}
