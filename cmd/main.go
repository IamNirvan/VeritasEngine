package main

import (
	"github.com/IamNirvan/veritasengine/internal/services/config"
	"github.com/IamNirvan/veritasengine/internal/services/data"
	log "github.com/sirupsen/logrus"
)

// This is the entrypoint of the service.
// All the dependencies will be obtained and the core application
// will be initialized and started here.
func main() {
	log.Debugf("service started")

	// Obtain the configuration object
	config := config.LoadConfig()

	// Obtain a connection with the data source
	data.NewDataSource(config)
}