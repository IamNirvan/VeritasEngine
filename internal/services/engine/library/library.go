package library

import (
	"sync"

	"github.com/IamNirvan/veritasengine/internal/enums"
	"github.com/IamNirvan/veritasengine/internal/services/config"
	rulesLoader "github.com/IamNirvan/veritasengine/internal/services/engine/loaders"
	"github.com/IamNirvan/veritasengine/internal/util/settings"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	log "github.com/sirupsen/logrus"
)

type LibraryManager interface {
	GetLibrary() *ast.KnowledgeLibrary
}

var (
	instance LibraryManager
	once     sync.Once
)

type LibraryManagerV1 struct {
	Config  *config.Config
	Library *ast.KnowledgeLibrary
}

func NewLibraryManager(cfg *config.Config) *LibraryManager {
	once.Do(func() {
		instance = &LibraryManagerV1{
			Config:  cfg,
			Library: nil,
		}
	})
	return &instance
}

// Loads the rules from the database and constructs a new library
// then returns it or nil if errors occurred
//
// Returns
// - *ast.KnowledgeLibrary if library was created or nil if errors occurred
func (lm *LibraryManagerV1) GetLibrary() *ast.KnowledgeLibrary {
	// Load the rules
	loader := rulesLoader.NewRulesLoader(lm.Config)
	rawRules, rulesErr := (*loader).LoadRules(enums.RULE_LOADING_FORMAT_STRING)
	if rulesErr != nil {
		return nil
	}

	// Make sure rules are in desired format
	rules, convertErr := rawRules.(string)
	if !convertErr {
		return nil
	}

	library := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(library)
	if buildErr := ruleBuilder.BuildRuleFromResource(settings.LIBRARY_KNOWLEDGE_BASE_NAME, settings.LIBRARY_VERSION, pkg.NewBytesResource([]byte(rules))); buildErr != nil {
		log.Warn("failed to construct library")
		return nil
	}
	return library
}
