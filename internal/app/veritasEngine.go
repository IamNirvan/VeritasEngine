package app

import (
	"context"

	"github.com/IamNirvan/veritasengine/internal/server"
	"github.com/IamNirvan/veritasengine/internal/services/config"
	"github.com/IamNirvan/veritasengine/internal/util"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type VeritasEngine struct {
	Config    *config.Config
	WebServer *server.WebServer
	Database  *gorm.DB
}

type VeritasEngineOpts struct {
	Config    *config.Config
	WebServer *server.WebServer
	Database  *gorm.DB
}

func NewVeritasEngine(opts *VeritasEngineOpts) *VeritasEngine {
	return &VeritasEngine{
		Config:    (*opts).Config,
		WebServer: (*opts).WebServer,
		Database:  (*opts).Database,
	}
}

func (ve *VeritasEngine) Start(ctx context.Context) error {
	errorGroup, ctx := errgroup.WithContext(ctx)

	// Create list of disposable resources
	// Add the disposables like the web server in here...
	disposables := []util.Disposable{ve.WebServer}

	// Run the web server in a goroutine
	errorGroup.Go(func() error {
		return ve.WebServer.Start()
	})

	// Create error message channel
	errChan := make(chan error)
	go func() {
		if err := errorGroup.Wait(); err != nil {
			errChan <- err
		}
	}()

	for {
		select {
		case <-ctx.Done():
			// Dispose all resources
			for index := range disposables {
				if err := disposables[index].Dispose(ctx); err != nil {
					return err
				}
			}
			return nil
		case err := <-errChan:
			return err
		}
	}
}
