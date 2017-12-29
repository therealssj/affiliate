package config

import (
	"fmt"

	"github.com/koding/multiconfig"
)

var initServerConfig = false
var serverConfig *ServerConfig

const BUY_COIN_UNIT_POWER int32 = 6

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
	ApiToken    string `teller`
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
	BuyerRate             float64   `default:"0.02"`
	LadderLine            []int     // default [0]
	PromoterRatio         []float64 //default [0.05]
	SuperiorPromoterRatio []float64 // default [0.03]
	SuperiorDiscount      float64   `default:"0.5"`
	MinSendAmount         int       `default:"1000000"`
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
	if len(daemonConfig.RewardConfig.LadderLine) == 0 {
		daemonConfig.RewardConfig.LadderLine = []int{0}
	}
	if len(daemonConfig.RewardConfig.PromoterRatio) == 0 {
		daemonConfig.RewardConfig.PromoterRatio = []float64{0.05}
	}
	if len(daemonConfig.RewardConfig.SuperiorPromoterRatio) == 0 {
		daemonConfig.RewardConfig.SuperiorPromoterRatio = []float64{0.03}
	}
	//	fmt.Printf("%+v\n", daemonConfig)
	if len(daemonConfig.RewardConfig.LadderLine) != len(daemonConfig.RewardConfig.PromoterRatio) || len(daemonConfig.RewardConfig.LadderLine) != len(daemonConfig.RewardConfig.SuperiorPromoterRatio) {
		panic("RewardConfig LadderLine, PromoterRatio, SuperiorPromoterRatio length not same")
	}
	initDaemonConfig = true
	return daemonConfig
}
