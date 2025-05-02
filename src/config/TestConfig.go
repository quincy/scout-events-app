package config

import "github.com/ridgelines/go-config"

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
		"database.hostname": "localhost",
		"database.username": "admin",
		"database.password": "admin",
		"database.port":     "26257",
		"database.name":     "scouting",
		"time.timezone":     "America/Boise",
	},
)
