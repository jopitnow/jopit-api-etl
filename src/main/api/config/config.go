package config

import (
	"fmt"
	"os"

	toolkitpkg "github.com/jopitnow/go-jopit-toolkit/gingonic/handlers"
	"github.com/jopitnow/go-jopit-toolkit/goutils/logger"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// test local env https://api-staging.jopit.com.ar

const (
	LogRatio     = 100
	LogBodyRatio = 100
)

var (
	InternalBaseItemsClient = "http://jopit-api-items:8080"
	InternalBaseShopsClient = "http://localhost:8081"
)

// Configuration structure
type Configuration struct {
	APIRestServerHost        string `mapstructure:"jopit_api_host"`
	APIRestServerPort        string `mapstructure:"jopit_api_port"`
	APIRestUsername          string `mapstructure:"jopit_api_username"`
	APIRestPassword          string `mapstructure:"jopit_api_password"`
	APIBaseEndpoint          string `mapstructure:"jopit_api_base_endpoint"`
	LoggingPath              string `mapstructure:"jopit_api_logpath"`
	LoggingFile              string `mapstructure:"jopit_api_logfile"`
	LoggingLevel             string `mapstructure:"jopit_api_loglevel"`
	MongoConnectionString    string `mapstructure:"MONGODB_CONN_STRING"`
	MercadolibreClientId     string `mapstructure:"MERCADOLIBRE_CLIENT_ID"`
	MercadolibreClientSecret string `mapstructure:"MERCADOLIBRE_CLIENT_SECRET"`
	AdminPassword            string
	AdminUsername            string
}

// ConfMap Config is package struct containing conf params
var ConfMap Configuration

func Load() {
	// Setting defaults if the config not read
	// API
	viper.SetDefault("jopit_api_host", "127.0.0.1")
	viper.SetDefault("jopit_api_port", ":8080")
	viper.SetDefault("jopit_api_username", "jopit")
	viper.SetDefault("jopit_api_password", "changeme")

	// LOG
	viper.SetDefault("jopit_api_logpath", "/var/log/jopit")
	viper.SetDefault("jopit_api_logfile", "jopit_api.log")
	viper.SetDefault("jopit_api_loglevel", "trace")

	// Read the config file
	viper.AutomaticEnv()

	err := viper.Unmarshal(&ConfMap)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %+v", err)
	}

	if os.Getenv("MONGODB_CONN_STRING") == "" {
		log.Fatal("MONGODB_CONN_STRING is empty)")
	}
	ConfMap.MongoConnectionString = os.Getenv("MONGODB_CONN_STRING")

	if os.Getenv("ADMIN_USERNAME") == "" {
		log.Fatal("ADMIN_USERNAME is empty)")
	}
	ConfMap.AdminUsername = os.Getenv("ADMIN_USERNAME")

	if os.Getenv("ADMIN_PASSWORD") == "" {
		log.Fatal("ADMIN_PASSWORD is empty)")
	}
	ConfMap.AdminPassword = os.Getenv("ADMIN_PASSWORD")

	if os.Getenv("MERCADOLIBRE_CLIENT_ID") == "" {
		log.Fatal("MERCADOLIBRE_CLIENT_ID is empty)")
	}
	ConfMap.MercadolibreClientId = os.Getenv("MERCADOLIBRE_CLIENT_ID")

	if os.Getenv("MERCADOLIBRE_CLIENT_SECRET") == "" {
		log.Fatal("MERCADOLIBRE_CLIENT_SECRET is empty)")
	}
	ConfMap.MercadolibreClientSecret = os.Getenv("MERCADOLIBRE_CLIENT_SECRET")

	toolkitpkg.ApiName = "items"

	os.Setenv("LOGGING_CONFIG_LEVEL", "0")
	os.Setenv("LOGGING_SAMPLING_LEVEL", "1")
	os.Setenv("LOGGING_LIMITER", "1500")

	logger.InitLoggerJopitConfig(toolkitpkg.ApiName)

	spew.Dump(ConfMap)

	fmt.Println("\n All good!!")
}
