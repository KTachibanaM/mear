package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"

	log "github.com/sirupsen/logrus"
)

var DockerContainerName = "mear-agent-testing"
var DockerNetworkName = "mear-network"

// Checking Docker exec's for 2 minutes max
var DockerExecCheckMaxAttempts = 60
var DockerExecCheckInterval = 2 * time.Second

type DevcontainerEngineProvisioner struct {
	container_id string
}

func NewDevcontainerEngineProvisioner() *DevcontainerEngineProvisioner {
	return &DevcontainerEngineProvisioner{
		container_id: "",
	}
}

func (p *DevcontainerEngineProvisioner) Provision(agent_binary_url, encoded_agent_args string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %v", err)
	}
	defer cli.Close()
	ctx := context.Background()

	log.Println("creating container...")
	container_create_resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "debian:bullseye",
		Cmd:   []string{"sleep", "infinity"},
	}, nil, nil, &v1.Platform{
		Architecture: "amd64",
	}, DockerContainerName)
	if err != nil {
		return fmt.Errorf("failed to create container: %v", err)
	}
	p.container_id = container_create_resp.ID

	log.Println("starting container...")
	err = cli.ContainerStart(ctx, p.container_id, types.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}

	log.Println("connecting container to network...")
	network_inspect_resp, err := cli.NetworkInspect(ctx, DockerNetworkName, types.NetworkInspectOptions{})
	if err != nil {
		return fmt.Errorf("failed to inspect network: %v", err)
	}
	err = cli.NetworkConnect(ctx, network_inspect_resp.ID, p.container_id, nil)
	if err != nil {
		return fmt.Errorf("failed to connect container to network: %v", err)
	}

	return nil
}

func (p *DevcontainerEngineProvisioner) Teardown() error {
	if p.container_id == "" {
		return fmt.Errorf("container was never provisioned")
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %v", err)
	}
	defer cli.Close()
	ctx := context.Background()

	log.Println("stopping container...")
	stop_timeout := 0
	err = cli.ContainerStop(ctx, p.container_id, container.StopOptions{Timeout: &stop_timeout})
	if err != nil {
		return err
	}

	log.Println("removing container...")
	err = cli.ContainerRemove(ctx, p.container_id, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		return err
	}

	return nil
}
