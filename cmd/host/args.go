package main

import "github.com/KTachibanaM/mear/internal/s3"

type HostArgs struct {
	InputFile         string       `json:"input_file"`
	DestinationTarget *s3.S3Target `json:"destination_target"`
	Stack             string       `json:"stack"`
	ExtraFfmpegArgs   []string     `json:"extra_ffmpeg_args"`
	DropletRam        int          `json:"droplet_ram"`
	DropletCpu        int          `json:"droplet_cpu"`
}

type HostResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func NewHostResult(success bool, message string) *HostResult {
	return &HostResult{
		Success: success,
		Message: message,
	}
}
