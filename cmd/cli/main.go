package main

import (
	"os"

	"github.com/KTachibanaM/mear/internal/host"
	"github.com/alexflint/go-arg"
	log "github.com/sirupsen/logrus"
)

type args struct {
	InputFile       string   `arg:"positional,required" help:"input file path"`
	OutputFile      string   `arg:"positional,required" help:"output file path"`
	Stack           string   `arg:"required" help:"stack for cloud resources. options are: dev, do"`
	RetainEngine    bool     `default:"false" help:"retain engine (VMs or containers) after running media encoder. helpful for debugging"`
	RetainBuckets   bool     `default:"false" help:"retain S3 buckets after running media encoder. helpful for debugging"`
	ExtraFfmpegArgs []string `help:"extra ffmpeg args to be passed for media encoding"`
}

func (args) Description() string {
	return "Self-hosted, on-demand cloud media encoding"
}

func main() {
	var args args
	p := arg.MustParse(&args)

	if _, err := os.Stat(args.InputFile); os.IsNotExist(err) {
		p.Fail("input file does not exist")
	}

	if args.Stack != "dev" && args.Stack != "do" {
		p.Fail("stack must be either dev or do")
	}

	err := host.Host(args.InputFile, args.OutputFile, args.RetainEngine, args.RetainBuckets)
	if err != nil {
		log.Fatalln(err)
	}
}
