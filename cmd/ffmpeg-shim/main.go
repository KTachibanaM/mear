package main

import (
	"os"

	"github.com/KTachibanaM/mear/internal/host"
	log "github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) != 4 {
		log.Fatalln("usage: mear <upload_filename> <save_to_filename> <deprovision_resources>")
	}

	err := host.Host(os.Args[1], os.Args[2], os.Args[3] == "true")
	if err != nil {
		log.Fatalln(err)
	}
}
