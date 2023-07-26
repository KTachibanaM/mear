package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

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

	if _, err := os.Stat(host_args.InputFile); os.IsNotExist(err) {
		fail("input file does not exist")
	}
	if host_args.Stack == "" {
		fail("stack must be specified")
	}
	if host_args.Stack != "dev" && host_args.Stack != "do" {
		fail("unknown stack name")
	}
	if host_args.Stack == "do" && (host_args.DropletRam == 0 || host_args.DropletCpu == 0) {
		fail("do ram and cpu must be specified")
	}

	err = host.Host(
		host_args.InputFile,
		host_args.DestinationTarget,
		host_args.Stack,
		false,
		false,
		host_args.ExtraFfmpegArgs,
		host_args.DropletRam,
		host_args.DropletCpu,
	)
	if err != nil {
		fail(err.Error())
	} else {
		success("successfully ran mear-host")
	}
}
