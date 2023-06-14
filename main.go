package main

import (
	"github.com/KTachibanaM/mear/agent"
	"github.com/KTachibanaM/mear/lib"
)

func main() {
	err := agent.Agent(
		lib.NewS3Target(
			"http://minio-source:9000",
			"us-east-1",
			"src",
			"MakeMine1948_256kb.rm",
			"minioadmin",
			"minioadmin",
			true,
		),
		lib.NewS3Target(
			"http://minio-destination:9000",
			"us-east-1",
			"dst",
			"output.mp4",
			"minioadmin",
			"minioadmin",
			true,
		),
	)
	if err != nil {
		panic(err)
	}
}
