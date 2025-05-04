package config

import (
	"errors"
	"fmt"
	"github.com/ridgelines/go-config"
	"log"
	"os"
	"strconv"
	"time"
)

type AppConfig interface {
	DatabaseHostname() string
	DatabaseName() string
	DatabasePassword() string
	DatabasePort() int
	DatabaseUsername() string
	Timezone() *time.Location
}

const dbHostname = "DB_HOSTNAME"
const dbName = "DB_NAME"
const dbPassword = "DB_PASSWORD"
const dbPort = "DB_PORT"
const dbUsername = "DB_USERNAME"
const timezone = "TIMEZONE"

var configKeys = map[string]string{
	dbHostname: "database.hostname",
	dbName:     "database.name",
	dbPassword: "database.password",
	dbPort:     "database.port",
	dbUsername: "database.username",
	timezone:   "time.timezone",
}

var defaults = map[string]string{
	dbHostname: "localhost",
	dbName:     "scouting",
	dbPassword: "admin",
	dbPort:     "26257",
	dbUsername: "admin",
	timezone:   "America/Boise",
}

type appConfig struct {
	cfg *config.Config
}

func ParseConfig() (AppConfig, error) {
	getenv, ok := os.LookupEnv("PROFILE")
	if !ok {
		getenv = "local"
	}
	var filePath string
	if getenv == "local" {
		filePath = "application-local.toml"
		log.Println("profile=local")
	} else if getenv == "prod" {
		filePath = "application-prod.toml"
		log.Println("profile=prod")
	} else {
		return nil, fmt.Errorf("unknown profile: %s", getenv)
	}

	configFile := config.NewTOMLFile(filePath)
	envConfig := config.NewEnvironment(configKeys)

	cfg := config.NewConfig([]config.Provider{configFile, envConfig})
	cfg.Validate = validate

	return &appConfig{cfg: cfg}, nil
}

func (a *appConfig) DatabaseHostname() string {
	value, err := a.cfg.StringOr(configKeys[dbHostname], defaults[dbHostname])
	if err != nil {
		log.Fatalf("Error getting database hostname: %v", err)
	}
	return value
}

func (a *appConfig) DatabaseName() string {
	value, err := a.cfg.StringOr(configKeys[dbName], defaults[dbName])
	if err != nil {
		log.Fatalf("Error getting database name: %v", err)
	}
	return value
}

func (a *appConfig) DatabasePassword() string {
	value, err := a.cfg.StringOr(configKeys[dbPassword], defaults[dbPassword])
	if err != nil {
		log.Fatalf("Error getting database password: %v", err)
	}
	return value
}

func (a *appConfig) DatabasePort() int {
	port, err := strconv.Atoi(defaults[dbPort])
	if err != nil {
		log.Fatalf("Error converting default port to int: %v", err)
	}
	value, err := a.cfg.IntOr(configKeys[dbPort], port)
	if err != nil {
		log.Fatalf("Error getting database port: %v", err)
	}
	return value
}

func (a *appConfig) DatabaseUsername() string {
	value, err := a.cfg.StringOr(configKeys[dbUsername], defaults[dbUsername])
	if err != nil {
		log.Fatalf("Error getting database username: %v", err)
	}
	return value
}

func (a *appConfig) Timezone() *time.Location {
	value, err := a.cfg.StringOr(configKeys[timezone], defaults[timezone])
	if err != nil {
		log.Fatalf("Error getting timezone: %v", err)
	}

	tz, err := time.LoadLocation(value)
	if err != nil {
		tz = time.UTC // fallback to UTC if timezone loading fails
		log.Printf("Error loading timezone %s: %v. Defaulting to UTC.", value, err)
	}

	return tz
}

func validate(settings map[string]string) error {
	var errs []error
	for _, v := range configKeys {
		if _, ok := settings[v]; !ok {
			errs = append(errs, fmt.Errorf("missing required config key: %s", v))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// TestProvider is a fake implementation of the config.Provider interface
type TestProvider struct {
	Values map[string]string
}

func NewTestConfig(testValues map[string]string) AppConfig {
	cfg := config.NewConfig([]config.Provider{TestProvider{Values: testValues}})
	return &appConfig{cfg: cfg}
}

func (t TestProvider) Load() (map[string]string, error) {
	return t.Values, nil
}

var TestConfig = NewTestConfig(
	map[string]string{
		dbHostname: "localhost",
		dbUsername: "admin",
		dbPassword: "admin",
		dbPort:     "26257",
		dbName:     "scouting",
		timezone:   "America/Boise",
	},
)
