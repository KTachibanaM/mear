package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/KTachibanaM/mear/internal/agent"
	"github.com/KTachibanaM/mear/internal/host"
	log "github.com/sirupsen/logrus"
)

func fail(msg string) {
	bytes, err := json.Marshal(NewHostResult(false, msg))
	if err != nil {
		fmt.Println("{\"success\": false, \"message\": \"failed to marshal json\"}")
	} else {
		fmt.Println(string(bytes))
	}
	os.Exit(1)
}

func failWithAgentFailure(agent_failure *agent.AgentFailure) {
	bytes, err := json.Marshal(NewHostResultWithAgentFailure(agent_failure))
	if err != nil {
		fmt.Println("{\"success\": false, \"message\": \"failed to marshal json\"}")
	} else {
		fmt.Println(string(bytes))
	}
	os.Exit(1)
}

func success(msg string) {
	bytes, err := json.Marshal(NewHostResult(true, msg))
	if err != nil {
		fmt.Println("{\"success\": true, \"message\": \"failed to marshal json\"}")
	} else {
		fmt.Println(string(bytes))
	}
	os.Exit(0)
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	reader := bufio.NewReader(os.Stdin)
	stdin, err := reader.ReadString('\n')
	if err != nil {
		fail("failed to read stdin")
	}

	var host_args HostArgs
	err = json.Unmarshal([]byte(stdin), &host_args)
	if err != nil {
		fail("failed to unmarshal json")
	}

	for _, host_job := range host_args.Jobs {
		if _, err := os.Stat(host_job.InputFile); os.IsNotExist(err) {
			fail(fmt.Sprintf("input file does not exist: %v", host_job.InputFile))
		}
	}
	if host_args.AgentExecutionTimeoutMinutes == 0 {
		fail("agent timeout must be specified")
	}
	if host_args.Stack == "" {
		fail("stack must be specified")
	}
	if host_args.Stack != "dev" && host_args.Stack != "do" {
		fail("unknown stack name")
	}
	if host_args.Stack == "do" {
		if host_args.DropletRam == 0 {
			fail("do ram must be specified")
		}
		if host_args.DropletCpu == 0 {
			fail("do cpu must be specified")
		}
		if host_args.DoAccessKeyId == "" {
			fail("do spaces access key id must be specified. go to https://cloud.digitalocean.com/account/api/spaces to generate one.")
		}
		if host_args.DoSecretAccessKey == "" {
			fail("do spaces secret access key must be specified. go to https://cloud.digitalocean.com/account/api/spaces to generate one.")
		}
		if host_args.DoToken == "" {
			fail("do token must be specified. go to https://cloud.digitalocean.com/account/api/tokens to generate one.")
		}
	}

	agent_failure, err := host.Host(
		host_args.Jobs,
		host_args.JobSharedDestinationS3Bucket,
		host_args.AgentExecutionTimeoutMinutes,
		host_args.Stack,
		false,
		false,
		host_args.DropletRam,
		host_args.DropletCpu,
		host_args.DoAccessKeyId,
		host_args.DoSecretAccessKey,
		host_args.DoToken,
	)
	if agent_failure != nil {
		failWithAgentFailure(agent_failure)
	} else if err != nil {
		fail(err.Error())
	} else {
		success("successfully ran mear-host")
	}
}
