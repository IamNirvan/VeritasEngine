package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/IamNirvan/veritasengine/internal/app"
	"github.com/IamNirvan/veritasengine/internal/handlers"
	"github.com/IamNirvan/veritasengine/internal/server"
	"github.com/IamNirvan/veritasengine/internal/services"
	"github.com/IamNirvan/veritasengine/internal/services/config"
	"github.com/IamNirvan/veritasengine/internal/services/data"
	log "github.com/sirupsen/logrus"
)

// This is the entrypoint of the service.
// All the dependencies will be obtained and the core application
// will be initialized and started here.
func main() {
	log.Debugf("service started")

	// Create context that listens for stop signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Obtain the configuration object
	config := config.LoadConfig()

	// Obtain a connection with the data source
	database := data.NewDataSource(config)

	// Initialize the services
	services := services.InitializeServices(&services.Options{
		Config:   config,
		Database: database,
	})

	// Initialize the handlers
	handlers := handlers.InitializeHandlers(&handlers.Options{
		Config:   config,
		Services: services,
	})

	// Initialize the web server
	webServer := server.NewWebServer(&server.WebServerOptions{
		Config:   config,
		Handlers: handlers,
	})

	// Create and start the application instance
	app := app.NewVeritasEngine(&app.VeritasEngineOpts{
		Config:    config,
		WebServer: webServer,
	})
	if err := app.Start(ctx); err != nil {
		log.Fatalf("failed to start application due to the following error: %s", err.Error())
	}
}
