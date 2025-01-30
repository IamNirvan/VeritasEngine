package ruleevaluation

import (
	"sync"

	"github.com/IamNirvan/veritasengine/internal/models/facts"
	"github.com/IamNirvan/veritasengine/internal/services"
	"github.com/IamNirvan/veritasengine/internal/services/config"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type RuleEvaluationHandler interface {
	EvaluateRule(ctx *gin.Context)
}

var (
	instance RuleEvaluationHandler
	once     sync.Once
)

type RuleEvaluationHandlerV1 struct {
	Config   *config.Config
	Services *services.Services
}

type RuleEvaluationOptions struct {
	Config   *config.Config
	Services *services.Services
}

func NewRuleEvaluationHandler(opts *RuleEvaluationOptions) *RuleEvaluationHandler {
	once.Do(func() {
		instance = &RuleEvaluationHandlerV1{
			Config:   (*opts).Config,
			Services: (*opts).Services,
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

	// Decode the fact
	var fact facts.GeneralInput
	if err := ctx.ShouldBindJSON(&fact); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Evaluate the rules using the decoded fact
	if response, err := (*handler.Services.RuleEvaluationService).EvaluateRule(&fact, ctx); err != nil {
		ctx.JSON(err.Status, gin.H{"error": err.Error})
		return
	} else {
		ctx.JSON(200, response)
	}
}
