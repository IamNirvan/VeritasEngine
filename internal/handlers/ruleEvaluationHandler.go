package handlers

import (
	"sync"

	"github.com/IamNirvan/veritasengine/internal/services/config"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type RuleEvaluationHandler interface {
	EvaluateRule(ctx *gin.Context)
}

type RuleEvaluationHandlerV1 struct {
	Config *config.Config
}

var (
	instance RuleEvaluationHandler
	once sync.Once
)

func NewRuleEvaluationHandler(cfg *config.Config) *RuleEvaluationHandler {
	once.Do(func() {
		instance = &RuleEvaluationHandlerV1 {
			Config: cfg,
		}
	})
	return &instance
}

// This function is the handler that evaluates the rules by
// using the rule evaluation service logic
//
// Parameters
// - ctx: a pointer to gin.Context
func (handler *RuleEvaluationHandlerV1) EvaluateRule(ctx *gin.Context) {
	log.Debug(" evaluating rules")

	// TODO: continue from here... call service logic here...
	ctx.JSON(200, "done");
} 