package host

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/KTachibanaM/mear/internal/agent"
	"github.com/KTachibanaM/mear/internal/s3"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/aws/aws-sdk-go/aws"
	aws_s3 "github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
)

var DockerContainerName = "mear-agent-testing"
var DockerNetworkName = "mear-network"

// 2 minutes
var DockerExecCheckMaxIntervals = 60
var DockerExecCheckInterval = 2 * time.Second

// 2 hours
var TailingMaxIntervals = 2880
var TailingInterval = 5 * time.Second

func Host() error {
	// 1. Get agent binary url
	log.Println("getting agent binary url...")
	agent_binary_url := "http://minio-agent-binary:9000/bin/mear-agent"

	// 2. Provision buckets
	log.Println("provisioning buckets...")
	s3_source := s3.NewS3Target(
		"http://minio-source:9000",
		"us-east-1",
		"src",
		"MakeMine1948_256kb.rm",
		"minioadmin",
		"minioadmin",
		true,
	)
	s3_destination := s3.NewS3Target(
		"http://minio-destination:9000",
		"us-east-1",
		"dst",
		"output.mp4",
		"minioadmin",
		"minioadmin",
		true,
	)
	s3_logs := s3.NewS3Target(
		"http://minio-destination:9000",
		"us-east-1",
		"dst",
		"agent.log",
		"minioadmin",
		"minioadmin",
		true,
	)

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
	next_log_position := 0
	for i := 0; i < TailingMaxIntervals; i++ {
		terminate := false
		bytes, err := s3.ReadS3Target(s3_logs)
		if err != nil {
			log.Warnf("failed to read logs, wait for next interval: %v", err)
		} else {
			logs := string(bytes)
			raw_log_lines := strings.Split(logs, "\n")
			log_lines := []string{}
			for _, raw_log_line := range raw_log_lines {
				trimmed_log_line := strings.Trim(raw_log_line, " ")
				if trimmed_log_line != "" {
					log_lines = append(log_lines, trimmed_log_line)
				}
			}
			if next_log_position >= len(log_lines) {
				log.Info("no new logs, wait for next interval")
			} else {
				for _, log_line := range log_lines[next_log_position:] {
					var structured_log map[string]interface{}
					err = json.Unmarshal([]byte(log_line), &structured_log)
					if err != nil {
						log.Warnf("failed to unmarshal log: %v", err)
					} else {
						msg, ok := structured_log["msg"].(string)
						if !ok {
							log.Warnf("failed to find msg in structured log")
						} else {
							log.Info(msg)
						}
						result, ok := structured_log["result"].(bool)
						if ok {
							if result {
								log.Info("agent ran successfully")
							} else {
								log.Errorln("agent run failed")
							}
							terminate = true
						}
					}
				}
				next_log_position = len(log_lines) + 1
			}
		}
		if terminate {
			break
		}
		time.Sleep(TailingInterval)
	}

	// 6. Deprovision buckets
	s3_sess, err := s3.CreateS3Session(s3_logs)
	if err != nil {
		return err
	}
	_, err = aws_s3.New(s3_sess).DeleteObject(
		&aws_s3.DeleteObjectInput{
			Bucket: aws.String(s3_logs.BucketName),
			Key:    aws.String(s3_logs.ObjectKey),
		},
	)
	if err != nil {
		return err
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
