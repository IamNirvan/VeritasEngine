package data

import (
	"sync"

	"github.com/IamNirvan/veritasengine/internal/services/config"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


var (
	instance *gorm.DB
	once     sync.Once
)

func NewDataSource(config *config.Config) *gorm.DB {
	once.Do(func() {
		// Get the connection string
		dsn := config.GetConnectionString()
		log.Tracef("Obtained connection string: %s", dsn)

		// Connect to the dataSource
		dataSource, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect to data source: %s", err.Error())
		}
		log.Debug("connected to data source")
		instance = dataSource
	})
	return instance
}