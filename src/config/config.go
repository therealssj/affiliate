package config

import (
	"fmt"

	"github.com/koding/multiconfig"
)

var initServerConfig = false
var serverConfig *ServerConfig

type ServerConfig struct {
	CoinName  string `default:"SPO TOKEN"`
	LogFolder string `default:"/tmp/affiliate/"`
	Db        Db
	Server    Server
	Teller    Teller
}

type Db struct {
	Host         string `default:"localhost"`
	Port         int    `default:"5432"`
	User         string `default:"lijt"`
	Password     string `default:"lijtlijt"`
	Name         string `default:"affiliate"`
	SslMode      string `default:"disable"`
	MaxOpenConns int    `default:"500"`
	MaxIdleConns int    `default:"50"`
}

type Server struct {
	Domain string `default:"localhost"`
	Port   int    `default:"6060"`
	Https  bool   `default:"false"`
}

type Teller struct {
	ContextPath string `default:"http://localhost:7071"`
}

func GetServerConfig() *ServerConfig {
	if initServerConfig {
		return serverConfig
	}
	m := multiconfig.NewWithPath("config.toml") // supports TOML, JSON and YAML
	serverConfig = new(ServerConfig)
	err := m.Load(serverConfig) // Check for error
	if err != nil {
		fmt.Printf("GetServerConfig Error: %s", err)
	}
	m.MustLoad(serverConfig) // Panic's if there is any error
	//	fmt.Printf("%+v\n", config)
	initServerConfig = true
	return serverConfig
}

var initDaemonConfig = false
var daemonConfig *DaemonConfig

type DaemonConfig struct {
	LogFolder    string `default:"/tmp/affiliate/"`
	Db           Db
	RewardConfig RewardConfig
	Teller       Teller
}

type RewardConfig struct {
	BuyerRate     float32
	PromoterRate  []float32
	MinSendAmount int
}

func GetDaemonConfig() *DaemonConfig {
	if initDaemonConfig {
		return daemonConfig
	}
	m := multiconfig.NewWithPath("config.toml") // supports TOML, JSON and YAML
	daemonConfig = new(DaemonConfig)
	err := m.Load(daemonConfig) // Check for error
	if err != nil {
		fmt.Printf("GetDaemonConfig Error: %s", err)
	}
	m.MustLoad(daemonConfig) // Panic's if there is any error
	//	fmt.Printf("%+v\n", daemonConfig)
	initDaemonConfig = true
	return daemonConfig
}
