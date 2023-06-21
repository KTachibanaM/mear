package bucket

import "github.com/KTachibanaM/mear/internal/s3"

var DevContainerSource = s3.NewS3Target(
	"http://minio-source:9000",
	"us-east-1",
	"src",
	"MakeMine1948_256kb.rm",
	"minioadmin",
	"minioadmin",
	true,
)
var DevContainerDestination = s3.NewS3Target(
	"http://minio-destination:9000",
	"us-east-1",
	"dst",
	"output.mp4",
	"minioadmin",
	"minioadmin",
	true,
)
var DevContainerLogs = s3.NewS3Target(
	"http://minio-destination:9000",
	"us-east-1",
	"dst",
	"agent.log",
	"minioadmin",
	"minioadmin",
	true,
)
