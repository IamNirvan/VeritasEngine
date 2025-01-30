package library

import (
	"sync"

	settings "github.com/IamNirvan/veritasengine/configs"
	"github.com/IamNirvan/veritasengine/internal/enums"
	"github.com/IamNirvan/veritasengine/internal/services/config"
	rulesLoader "github.com/IamNirvan/veritasengine/internal/services/engine/loaders"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
	DB      *gorm.DB
	Library *ast.KnowledgeLibrary
}

func NewLibraryManager(cfg *config.Config, db *gorm.DB) *LibraryManager {
	once.Do(func() {
		instance = &LibraryManagerV1{
			Config:  cfg,
			DB:      db,
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
	loader := rulesLoader.NewRulesLoader(lm.Config, lm.DB)
	rawRules, rulesErr := (*loader).LoadRules(enums.RULE_LOADING_FORMAT_STRING)
	if rulesErr != nil {
		return nil
	}

	// Make sure rules are in desired format
	// Make sure rules are in desired format
	var rules string
	switch v := rawRules.(type) {
	case string:
		rules = v
	case []byte:
		rules = string(v)
	case *string:
		if v != nil {
			rules = *v
		} else {
			log.Warn("failed to convert rules to string, *string is nil")
			return nil
		}
	default:
		log.Warnf("failed to convert rules to string, unexpected type: %T", rawRules)
		return nil
	}

	log.Debugf("loaded rules after converting: %s", rules)

	log.Debugf("loaded rules after converting: %s", rules)

	library := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(library)
	if buildErr := ruleBuilder.BuildRuleFromResource(settings.LIBRARY_KNOWLEDGE_BASE_NAME, settings.LIBRARY_VERSION, pkg.NewBytesResource([]byte(rules))); buildErr != nil {
		log.Warn("failed to construct library")
		return nil
	}

	log.Debugf("constructed library: %v", library)
	return library
}
