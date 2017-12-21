package config

import (
	"fmt"

	"github.com/koding/multiconfig"
)

var initConfig = false
var config *Config

type Config struct {
	Db     Db
	Server Server
}

type Db struct {
	Host     string `default:"localhost"`
	Port     int    `default:"5432"`
	User     string `default:"lijt"`
	Password string `default:"lijtlijt"`
	Name     string `default:"affiliate"`
	SslMode  string `default:"disable"`
}

type Server struct {
	Domain string `default:"localhost"`
	Port   int    `default:"6060"`
	Https  bool   `default:false`
}

func GetConfig() *Config {
	if initConfig {
		return config
	}
	m := multiconfig.NewWithPath("config.toml") // supports TOML, JSON and YAML
	// Get an empty struct for your configuration
	config := new(Config)
	// Populated the serverConf struct
	err := m.Load(config) // Check for error
	if err != nil {
		fmt.Printf("GetConfig Error: %s", err)
	}
	m.MustLoad(config) // Panic's if there is any error
	//	fmt.Printf("%+v\n", config)
	return config
}
