package main

import (
	"github.com/KTachibanaM/mear/internal/host"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := host.Host()
	if err != nil {
		log.Fatalf("failed to run host: %v", err)
	}
}
