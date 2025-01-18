package ruleevaluation

import (
	"context"

	"github.com/IamNirvan/veritasengine/internal/errors"
	"github.com/IamNirvan/veritasengine/internal/services/config"
	"github.com/IamNirvan/veritasengine/internal/services/engine"
	"github.com/IamNirvan/veritasengine/internal/services/engine/library"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RuleEvaluationService interface {
	EvaluateRule(interface{}, context.Context) (*interface{}, *errors.ServiceError)
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

func (res *RuleEvaluationServiceV1) EvaluateRule(fact interface{}, ctx context.Context) (*interface{}, *errors.ServiceError) {
	// Get the updated library
	libraryManager := library.NewLibraryManager(res.Config)
	lib := (*libraryManager).GetLibrary()

	// Fetch knowledge base
	knowledgeBase, knowledgeBaseErr := lib.NewKnowledgeBaseInstance(library.KNOWLEDGE_BASE_NAME, library.VERSION)
	if knowledgeBaseErr != nil {
		log.Errorf("error when obtaining instance of KnowledgeBase: %v", knowledgeBaseErr.Error())
		return nil, &errors.ServiceError{Error: knowledgeBaseErr.Error(), Status: 500}
	}

	// Get the rule engine
	engine := engine.NewRuleEngine()

	// Create a data context
	dataCtx := ast.NewDataContext()
	if dataCtxErr := dataCtx.Add("DDF", fact); dataCtxErr != nil {
		log.Errorf("error when adding fact to data context: %v", dataCtxErr.Error())
		return nil, &errors.ServiceError{Error: dataCtxErr.Error(), Status: 500}
	}

	// Evaluate rule(s)
	if engineErr := engine.Execute(dataCtx, knowledgeBase); engineErr != nil {
		log.Errorf("error when evaluating rule: %v", engineErr.Error())
		return nil, &errors.ServiceError{Error: engineErr.Error(), Status: 500}
	}

	// TODO: finalize the response
	return nil, nil
}
