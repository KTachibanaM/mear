package host

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/KTachibanaM/mear/internal/agent"
	"github.com/KTachibanaM/mear/internal/bucket"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"

	log "github.com/sirupsen/logrus"
)

var DockerContainerName = "mear-agent-testing"
var DockerNetworkName = "mear-network"

// 2 minutes
var DockerExecCheckMaxIntervals = 60
var DockerExecCheckInterval = 2 * time.Second

func Host() error {
	// 1. Get agent binary url
	log.Println("getting agent binary url...")
	agent_binary_url := "http://minio-agent-binary:9000/bin/mear-agent"

	// 2. Provision buckets
	log.Println("provisioning buckets...")
	bucket_provisioner := bucket.NewDevcontainerBucketProvisioner()
	s3_source, s3_destination, s3_logs, err := bucket_provisioner.Provision()

	if err != nil {
		log.Fatalf("failed to provision buckets: %v", err)
	}

	// 3. Gather agent args
	log.Println("gathering agent args...")
	agent_args := agent.NewAgentArgs(
		s3_source,
		s3_destination,
		s3_logs,
		[]string{},
	)
	agent_args_json, err := json.MarshalIndent(agent_args, "", "")
	if err != nil {
		return err
	}
	encoded := base64.StdEncoding.EncodeToString(agent_args_json)

	// 4. Provision engine
	log.Println("provisioning engine...")
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()

	ctx := context.Background()
	println("create")
	container_create_resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "debian:bullseye",
		Cmd:   []string{"sleep", "infinity"},
	}, nil, nil, &v1.Platform{
		Architecture: "amd64",
	}, DockerContainerName)
	if err != nil {
		return err
	}
	println("start")
	err = cli.ContainerStart(ctx, container_create_resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	println("nw inspect")
	network_inspect_resp, err := cli.NetworkInspect(ctx, DockerNetworkName, types.NetworkInspectOptions{})
	if err != nil {
		return err
	}
	println("nw connect")
	err = cli.NetworkConnect(ctx, network_inspect_resp.ID, container_create_resp.ID, nil)
	if err != nil {
		return err
	}
	commands := []string{
		"apt update",
		"apt install -y curl",
		"curl -sL " + agent_binary_url + " -o /root/mear-agent",
		"chmod +x /root/mear-agent",
		"/root/mear-agent " + encoded,
	}
	for i := 0; i < len(commands); i++ {
		command := commands[i]
		fmt.Printf("running %v\n", command)
		exec_create_resp, err := cli.ContainerExecCreate(context.Background(), container_create_resp.ID, types.ExecConfig{
			Cmd: []string{"sh", "-c", command},
		})
		if err != nil {
			return err
		}

		err = cli.ContainerExecStart(context.Background(), exec_create_resp.ID, types.ExecStartCheck{})
		if err != nil {
			return err
		}
		if i != len(commands)-1 {
			for j := 0; j < DockerExecCheckMaxIntervals; j++ {
				execInspectResp, err := cli.ContainerExecInspect(ctx, exec_create_resp.ID)
				if err != nil {
					return err
				}

				if !execInspectResp.Running {
					break
				}
				time.Sleep(DockerExecCheckInterval)
			}
		}
	}

	// 5. Tail for logs and result
	log.Println("tailing for logs and result...")
	result := TailS3Logs(s3_logs)
	if result {
		log.Println("agent run succeeded")
	} else {
		log.Println("agent run failed")
	}

	// 6. Deprovision buckets
	err = bucket_provisioner.Teardown()
	if err != nil {
		log.Fatalf("failed to teardown buckets: %v", err)
	}

	// 7. Deprovision engine
	log.Println("deprovisioning engine...")
	err = cli.ContainerStop(ctx, container_create_resp.ID, container.StopOptions{})
	if err != nil {
		return err
	}
	err = cli.ContainerRemove(ctx, container_create_resp.ID, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}

	return nil
}
