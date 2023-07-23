package main

import (
	"fmt"
	"os"

	"github.com/KTachibanaM/mear/internal/host"
)

func usage() {
	fmt.Println("Usage: mear -i <input file> --mear-stack <stack name> [--mear-retain-engine] [--mear-retain-buckets] [extra ffmpeg args] <output file>")
	os.Exit(1)
}

func fail(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func main() {
	var input string
	var output string
	var stack string
	retainEngine := false
	retainBuckets := false
	var extraFfmpegArgs []string

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
		} else if arg == "--mear-stack" {
			stack = args[i+1]
			i++
		} else if arg == "--mear-retain-engine" {
			retainEngine = true
		} else if arg == "--mear-retain-buckets" {
			retainBuckets = true
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
	if stack == "" {
		usage()
	}
	if stack != "dev" && stack != "do" {
		fail("unknown stack name")
	}

	err := host.Host(input, output, stack, retainEngine, retainBuckets, extraFfmpegArgs)
	if err != nil {
		fail(err.Error())
	}
}
