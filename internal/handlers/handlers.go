package handlers

import (
	"sync"

	ruleevaluation "github.com/IamNirvan/veritasengine/internal/handlers/ruleEvaluation"
	"github.com/IamNirvan/veritasengine/internal/services"
	"github.com/IamNirvan/veritasengine/internal/services/config"
)

type Handlers struct {
	RuleEvaluationHandler *ruleevaluation.RuleEvaluationHandler
}

var (
	instance *Handlers
	once     sync.Once
)

type Options struct {
	Config   *config.Config
	Services *services.Services
}

func InitializeHandlers(opts *Options) *Handlers {
	once.Do(func() {
		instance = &Handlers{
			RuleEvaluationHandler: ruleevaluation.NewRuleEvaluationHandler(&ruleevaluation.RuleEvaluationOptions{
				Config:   (*opts).Config,
				Services: (*opts).Services,
			}),
		}
	})
	return instance
}
