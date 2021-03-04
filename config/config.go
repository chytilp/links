package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/chytilp/links/logging"
)

// ConfigFile is the default environment variable the points to the config file.
const ConfigFile = "config.toml"

// App is the application config.
var App *Config

// Config defines structs for the config file.
type Config struct {
	// Database connection string components.
	Database DbConfig `toml:"database"`
}

// DbConfig is database connection string components struct.
type DbConfig struct {
	Address  string
	Port     int
	Database string
	User     string
	Password string
}

// GetConnectionString func formats Database string components into connection string.
func (d *DbConfig) GetConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", d.User, d.Password, d.Address, d.Port, d.Database)
}

// init func loads and returns the config loaded from environment.
func init() {
	App = &Config{}
	if _, err := toml.DecodeFile(ConfigFile, App); err != nil {
		logging.L.Error("Error from read config.toml file. err: %s", err)
		fmt.Println(err)
		return
	}
}
