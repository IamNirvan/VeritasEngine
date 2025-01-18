package services

import (
	"sync"

	"github.com/IamNirvan/veritasengine/internal/services/config"
	ruleevaluation "github.com/IamNirvan/veritasengine/internal/services/ruleEvaluation"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Options struct {
	Config   *config.Config
	Database *gorm.DB
}

type Services struct {
	RuleEvaluationService *ruleevaluation.RuleEvaluationService
}

var (
	instance *Services
	once     sync.Once
)

func InitializeServices(opts *Options) *Services {
	once.Do(func() {
		instance = &Services{
			RuleEvaluationService: ruleevaluation.NewRuleEvaluationServiceV1(&ruleevaluation.RuleEvaluationOptions{
				Config:   (*opts).Config,
				Database: (*opts).Database,
			}),
		}
		log.Trace("initialized services")
	})
	return instance
}
