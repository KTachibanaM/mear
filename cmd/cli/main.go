package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/KTachibanaM/mear/internal/cli"
)

func usage() {
	fmt.Println("Usage: mear -i <input file> --mear-stack <stack name> --mear-agent-timeout <timeout minutes> [--mear-retain-engine] [--mear-retain-buckets] [--mear-do-ram] <ram in gb> [--mear-do-cpu] <cpu core count> [extra ffmpeg args] <output file>")
	fmt.Println("       -i <input file>:                input file path")
	fmt.Println("       --mear-stack <stack name>:      cloud stack to use for media encoding")
	fmt.Println("                                       options are 'do' (DigitalOcean) and 'dev' (development in devcontainer)")
	fmt.Println("       --mear-agent-timeout <minutes>: agent execution timeout in minutes")
	fmt.Println("       --mear-retain-engine:           retain the engine (VPS or container) after media encoding")
	fmt.Println("                                       default is false")
	fmt.Println("       --mear-retain-buckets:          retain the S3 buckets after media encoding")
	fmt.Println("                                       default is false")
	fmt.Println("       --mear-do-ram <ram in gb>:      DigitalOcean droplet ram in gb")
	fmt.Println("                                       mandatory if --mear-stack is 'do'")
	fmt.Println("                                       supported combinations with --mear-do-cpu are 1gb/1cpu, 2gb/1cpu, 2gb/2cpu, 4gb/2cpu, 8gb/4cpu and 16gb/8cpu")
	fmt.Println("       --mear-do-cpu <cpu core count>: DigitalOcean droplet cpu core count")
	fmt.Println("                                       mandatory if --mear-stack is 'do'")
	fmt.Println("                                       supported combinations with --mear-do-ram are 1gb/1cpu, 2gb/1cpu, 2gb/2cpu, 4gb/2cpu, 8gb/4cpu and 16gb/8cpu")
	fmt.Println("       [extra ffmpeg args]:            extra args to pass into ffmpeg for media encoding")
	fmt.Println("       <output file>:                  output file path")
	os.Exit(1)
}

func fail(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func main() {
	var input string
	var output string
	var agentExecutionTimeoutMinutes int
	var stack string
	retainEngine := false
	retainBuckets := false
	var doRam int
	var doCpu int
	var extraFfmpegArgs []string
	var doAccessKeyId string
	var doSecretAccessKey string
	var doToken string

	args := os.Args[1:]
	if len(args) == 0 {
		usage()
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if i == len(args)-1 {
			output = arg
		} else if arg == "-i" {
			input = args[i+1]
			i++
		} else if arg == "--mear-agent-timeout" {
			agentExecutionTimeoutMinutesStr := args[i+1]
			agentExecutionTimeoutMinutes, _ = strconv.Atoi(agentExecutionTimeoutMinutesStr)
			i++
		} else if arg == "--mear-stack" {
			stack = args[i+1]
			i++
		} else if arg == "--mear-retain-engine" {
			retainEngine = true
		} else if arg == "--mear-retain-buckets" {
			retainBuckets = true
		} else if arg == "--mear-do-ram" {
			doRamStr := args[i+1]
			doRam, _ = strconv.Atoi(doRamStr)
			i++
		} else if arg == "--mear-do-cpu" {
			doCpuStr := args[i+1]
			doCpu, _ = strconv.Atoi(doCpuStr)
			i++
		} else {
			extraFfmpegArgs = append(extraFfmpegArgs, arg)
		}
	}

	if _, err := os.Stat(input); os.IsNotExist(err) {
		fail("input file does not exist")
	}
	if output == "" {
		usage()
	}
	if agentExecutionTimeoutMinutes == 0 {
		fail("agent timeout must be specified")
	}
	if stack == "" {
		usage()
	}
	if stack != "dev" && stack != "do" {
		fail("unknown stack name")
	}
	if stack == "do" {
		if doRam == 0 {
			fail("do ram must be specified")
		}
		if doCpu == 0 {
			fail("do cpu must be specified")
		}
		doAccessKeyId = get_config(
			"do",
			"access_key_id",
			"DigitalOcean Spaces access key id. Go to https://cloud.digitalocean.com/account/api/spaces to generate one",
		)
		if doAccessKeyId == "" {
			fail("do spaces access key id must be specified")
		}
		doSecretAccessKey = get_config(
			"do",
			"secret_access_key",
			"DigitalOcean Spaces secret access key. Go to https://cloud.digitalocean.com/account/api/spaces to generate one",
		)
		if doSecretAccessKey == "" {
			fail("do spaces secret access key must be specified")
		}
		doToken = get_config(
			"do",
			"token",
			"DigitalOcean API token. Go to https://cloud.digitalocean.com/account/api/tokens to generate one",
		)
		if doToken == "" {
			fail("do token must be specified")
		}
	}

	err := cli.Cli(
		input,
		output,
		extraFfmpegArgs,
		agentExecutionTimeoutMinutes,
		stack,
		retainEngine,
		retainBuckets,
		doRam,
		doCpu,
		doAccessKeyId,
		doSecretAccessKey,
		doToken,
	)
	if err != nil {
		fail(err.Error())
	}
}
