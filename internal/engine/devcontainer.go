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
var DockerImage = "debian:bullseye"

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

func (p *DevcontainerEngineProvisioner) Provision(agent_binary_url string, ssh_public_key []byte) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", fmt.Errorf("failed to create docker client: %v", err)
	}
	defer cli.Close()
	ctx := context.Background()

	log.Println("creating container...")
	container_create_resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: DockerImage,
		Cmd:   []string{"sleep", "infinity"},
	}, nil, nil, &v1.Platform{
		Architecture: "amd64",
	}, DockerContainerName)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %v", err)
	}
	p.container_id = container_create_resp.ID

	log.Println("starting container...")
	err = cli.ContainerStart(ctx, p.container_id, types.ContainerStartOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to start container: %v", err)
	}

	log.Println("connecting container to network...")
	network_inspect_resp, err := cli.NetworkInspect(ctx, DockerNetworkName, types.NetworkInspectOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to inspect network: %v", err)
	}
	err = cli.NetworkConnect(ctx, network_inspect_resp.ID, p.container_id, nil)
	if err != nil {
		return "", fmt.Errorf("failed to connect container to network: %v", err)
	}

	commands := []string{
		"apt update",
		"apt install -y openssh-server",
		"mkdir -p /root/.ssh",
		fmt.Sprintf("echo \"%v\" > /root/.ssh/authorized_keys", ssh_public_key),
		"service ssh start",
	}
	for i := 0; i < len(commands); i++ {
		command := commands[i]
		log.Infof("executing in container the command '%v'\n", command)

		exec_create_resp, err := cli.ContainerExecCreate(context.Background(), p.container_id, types.ExecConfig{
			Cmd: []string{"sh", "-c", command},
		})
		if err != nil {
			return "", fmt.Errorf("failed to create exec: %v", err)

		}

		err = cli.ContainerExecStart(context.Background(), exec_create_resp.ID, types.ExecStartCheck{})
		if err != nil {
			return "", fmt.Errorf("failed to start exec: %v", err)

		}
		for j := 0; j < DockerExecCheckMaxAttempts; j++ {
			exec_inspect_resp, err := cli.ContainerExecInspect(ctx, exec_create_resp.ID)
			if err != nil {
				return "", fmt.Errorf("failed to inspect exec: %v", err)

			}
			if !exec_inspect_resp.Running {
				break
			}
			time.Sleep(DockerExecCheckInterval)
		}
	}

	log.Println("getting container's IP address...")
	container, err := cli.ContainerInspect(context.Background(), p.container_id)
	if err != nil {
		return "", fmt.Errorf("failed to inspect container: %v", err)
	}

	return container.NetworkSettings.IPAddress, nil
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
