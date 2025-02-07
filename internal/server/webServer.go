package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	settings "github.com/IamNirvan/veritasengine/configs"
	"github.com/IamNirvan/veritasengine/internal/handlers"
	"github.com/IamNirvan/veritasengine/internal/services/config"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type WebServer struct {
	Config   *config.Config
	Handlers *handlers.Handlers
	Server   *http.Server
}

var (
	instance *WebServer
	once     sync.Once
)

type WebServerOptions struct {
	Config   *config.Config
	Handlers *handlers.Handlers
}

func NewWebServer(opts *WebServerOptions) *WebServer {
	once.Do(func() {
		instance = &WebServer{
			Config:   (*opts).Config,
			Handlers: (*opts).Handlers,
		}
		log.Trace("initialized web server")
	})
	return instance
}

func (ws *WebServer) Start() error {
	// Run Gin in production mode if the config mode is set to prod
	if (*ws.Config).Mode == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Create a router group under the engine's version 1 base url
	version1Base := r.Group(settings.SERVER_BASE_URL_V1)
	version1Base.POST("/evaluate/rule", (*ws.Handlers.RuleEvaluationHandler).EvaluateRule)

	ws.Server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", ws.Config.WebServer.Host, ws.Config.WebServer.Port),
		Handler: r,
	}

	log.Debug("started web server")
	return ws.Server.ListenAndServe()
}

func (ws *WebServer) Dispose(ctx context.Context) error {
	err := ws.Server.Shutdown(ctx)
	log.Debug("stopped web server")
	return err
}
