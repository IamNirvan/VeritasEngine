package rulesLoader

import (
	"fmt"
	"strings"
	"sync"

	"github.com/IamNirvan/veritasengine/internal/services/config"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RulesLoader interface {
	LoadRules(int) (interface{}, error)
}

type RulesLoaderV1 struct {
	Config *config.Config
	DB     *gorm.DB
}

var (
	instance RulesLoader
	once     sync.Once
)

func NewRulesLoader(cfg *config.Config, db *gorm.DB) *RulesLoader {
	once.Do(func() {
		instance = &RulesLoaderV1{
			Config: cfg,
			DB:     db,
		}
	})
	return &instance
}

// This loads the rules according to the specified format
// The formats can be found in the ruleLoadingFormats.go file
//
// Parameters
// - format: an enum value
//
// Returns
// - interface{} containing the rules in the specified format
func (rl *RulesLoaderV1) LoadRules(format int) (interface{}, error) {
	var rules interface{}
	var err error

	if format == 0 {
		rules, err = rl.loadRulesAsString()
	} else {
		err = fmt.Errorf("invalid rule loading format: %s", format)
	}

	if err != nil {
		return nil, err
	}
	return rules, nil
}

func (rl *RulesLoaderV1) loadRulesAsString() (*string, error) {
	loadedRules, err := rl.loadRulesFromDB()
	if err != nil {
		return nil, err
	}

	// Build the rules as a string
	var sb strings.Builder
	for _, rule := range *loadedRules {
		sb.WriteString(rule)
		sb.WriteString("\n")
	}

	result := sb.String()
	return &result, nil
}

func (rl *RulesLoaderV1) loadRulesFromDB() (*[]string, error) {
	var loadedRules []string
	if err := rl.DB.Table("rules").Select("rule").Find(&loadedRules).Error; err != nil {
		return nil, fmt.Errorf("failed to load rules from database: %v", err)
	}
	log.Tracef("loaded %d rule(s)", len(loadedRules))
	return &loadedRules, nil
}
