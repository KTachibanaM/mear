package main

import (
	"os"

	"github.com/KTachibanaM/mear/internal/host"
	log "github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) != 5 {
		log.Fatalln("usage: mear <upload_filename> <save_to_filename> <skip_deprovision_engine> <skip_deprovision_buckets>")
	}

	err := host.Host(os.Args[1], os.Args[2], os.Args[3] == "true", os.Args[4] == "true")
	if err != nil {
		log.Fatalln(err)
	}
}
