package ruleevaluation

import (
	"context"

	settings "github.com/IamNirvan/veritasengine/configs"
	"github.com/IamNirvan/veritasengine/internal/errors"
	"github.com/IamNirvan/veritasengine/internal/models/facts"
	"github.com/IamNirvan/veritasengine/internal/services/config"
	"github.com/IamNirvan/veritasengine/internal/services/engine"
	"github.com/IamNirvan/veritasengine/internal/services/engine/library"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RuleEvaluationService interface {
	EvaluateRule(*facts.GeneralInput, context.Context) (*[]interface{}, *errors.ServiceError)
}

type RuleEvaluationServiceV1 struct {
	Config   *config.Config
	Database *gorm.DB
}

type RuleEvaluationOptions struct {
	Config   *config.Config
	Database *gorm.DB
}

func NewRuleEvaluationServiceV1(opts *RuleEvaluationOptions) *RuleEvaluationService {
	var service RuleEvaluationService = &RuleEvaluationServiceV1{
		Config:   (*opts).Config,
		Database: (*opts).Database,
	}
	return &service
}

func (res *RuleEvaluationServiceV1) EvaluateRule(fact *facts.GeneralInput, ctx context.Context) (*[]interface{}, *errors.ServiceError) {
	// Get the updated library
	libraryManager := library.NewLibraryManager(res.Config, res.Database)
	lib := (*libraryManager).GetLibrary()
	if lib == nil {
		log.Error("error when obtaining library")
		return nil, &errors.ServiceError{Error: "error when obtaining library", Status: 500}
	}

	// Fetch knowledge base
	knowledgeBase, knowledgeBaseErr := lib.NewKnowledgeBaseInstance(settings.LIBRARY_KNOWLEDGE_BASE_NAME, settings.LIBRARY_VERSION)
	if knowledgeBaseErr != nil {
		log.Errorf("error when obtaining instance of KnowledgeBase: %v", knowledgeBaseErr.Error())
		return nil, &errors.ServiceError{Error: knowledgeBaseErr.Error(), Status: 500}
	}

	// Get the rule engine
	engine := engine.NewRuleEngine()

	// Create a data context
	dataCtx := ast.NewDataContext()
	if dataCtxErr := dataCtx.Add(settings.FACT_ALIAS, fact); dataCtxErr != nil {
		log.Errorf("error when adding fact to data context: %v", dataCtxErr.Error())
		return nil, &errors.ServiceError{Error: dataCtxErr.Error(), Status: 500}
	}

	// Evaluate rule(s)
	if engineErr := engine.Execute(dataCtx, knowledgeBase); engineErr != nil {
		log.Errorf("error when evaluating rule: %v", engineErr.Error())
		return nil, &errors.ServiceError{Error: engineErr.Error(), Status: 500}
	}

	// Fetch the response from the fact
	response := &fact.Response
	log.Tracef("evaluated rules response: %+v", response)

	return *response, nil
}
