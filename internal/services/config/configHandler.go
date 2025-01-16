package config

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/IamNirvan/veritasengine/internal/models"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Log       models.Log       `mapstructure:"log"`
	WebServer models.WebServer `mapstructure:"web_server"`
	Database  models.Database  `mapstructure:"database"`
}

var (
	once sync.Once
	instance *Config
)


func newConfig() *Config {
	return &Config{
		// These will be overriden by config file or environment variables
		Log: models.Log{
			Level:   "debug",
			Methods: true,
		},
		WebServer: models.WebServer{
			Host:    "localhost",
			Port:    8080,
			Timeout: 15,
		},
		Database: models.Database{
			User:     "postgres",
			Password: "postgres",
			Host:     "localhost",
			Port:     5432,
			Dbname:   "postgres",
			Sslmode:  "disable",
		},
	}
}

// Loads the configration values from the config file or environment
// variables.
//
// Returns
// - A pointer to a struct with configuration attibutes
func LoadConfig() *Config {
	once.Do(func () {
		vpr := viper.NewWithOptions(viper.EnvKeyReplacer(strings.NewReplacer(".", "_")))

		vpr.SetConfigName("config")
		vpr.SetConfigType("yaml")
	
		vpr.AddConfigPath("./configs")
		vpr.AddConfigPath(".")
	
		exe, err := os.Executable()
		if err != nil {
			log.Fatalf("cannot access executable directory, error : %v", err)
		}
		vpr.AddConfigPath(path.Dir(exe))
	
		if err := vpr.ReadInConfig(); err != nil {
			log.Fatalf("config file reading error : %v", err)
		}
	
		vpr.SetEnvPrefix("VERITAS_ENGINE")
		vpr.AutomaticEnv()
	
		config := newConfig()
		if err := vpr.Unmarshal(config, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(mapstructure.TextUnmarshallerHookFunc()))); err != nil {
			log.Fatalf("failed to decode configuration to struct with error: %v", err)
		}
	
		var level log.Level
		if err := level.UnmarshalText([]byte(config.Log.Level)); err != nil {
			log.Fatal("log level not valid, error:", err)
		}
		log.SetLevel(level)
		log.SetReportCaller(config.Log.Methods)

		instance = config
		log.Debug("obtained configuration")
	})
	return instance
}

func (c *Config) GetConnectionString() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Dbname,
		c.Database.Sslmode,
	)
}
