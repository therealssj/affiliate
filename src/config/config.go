package config

import (
	"fmt"

	"github.com/koding/multiconfig"
)

var initServerConfig = false
var serverConfig *ServerConfig

type ServerConfig struct {
	CoinName      string `default:"SPO TOKEN"`
	LogFolder     string `default:"/tmp/affiliate/"`
	Db            Db
	Server        Server
	Teller        Teller
	CoinUnitPower int `default:"6"`
}

type Db struct {
	Host          string `default:"localhost"`
	Port          int    `default:"5432"`
	User          string `default:"lijt"`
	Password      string `default:"lijtlijt"`
	Name          string `default:"affiliate"`
	SslMode       string `default:"disable"`
	MaxOpenConns  int    `default:"500"`
	MaxIdleConns  int    `default:"50"`
	ChecksumToken string `default:"test-checksum-token"` // for testing convenience, must be specified other string in config.toml
}

type Server struct {
	Domain     string `default:"localhost"`
	Port       int    `default:"80"`
	ListenIp   string `default:"127.0.0.1"`
	ListenPort int    `default:"6060"`
	Https      bool   `default:"false"`
}

type Teller struct {
	ContextPath string `default:"http://localhost:7071"`
	ApiToken    string
	Debug       bool `default:"false"`
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

var initApiForTellerConfig = false
var apiForTellerConfig *ApiForTellerConfig

type ApiForTellerConfig struct {
	LogFolder    string `default:"/tmp/affiliate/"`
	Db           Db
	RewardConfig RewardConfig
	ListenIp     string `default:"127.0.0.1"`
	ListenPort   int    `default:"6010"`
	AuthToken    string
	AuthValidSec int  `default:"15"`
	Debug        bool `default:"false"`
}

type RewardConfig struct {
	BuyerRatio            float64   `default:"0.02"`
	LadderLine            []int     // default [0]
	PromoterRatio         []float64 //default [0.05]
	SuperiorPromoterRatio []float64 // default [0.03]
	SuperiorDiscount      float64   `default:"0.5"`
	MinSendAmount         int       `default:"1000000"`
}

func GetApiForTellerConfig() *ApiForTellerConfig {
	if initApiForTellerConfig {
		return apiForTellerConfig
	}
	m := multiconfig.NewWithPath("config.toml") // supports TOML, JSON and YAML
	apiForTellerConfig = new(ApiForTellerConfig)
	err := m.Load(apiForTellerConfig) // Check for error
	if err != nil {
		fmt.Printf("GetApiForTellerConfig Error: %s", err)
	}
	m.MustLoad(apiForTellerConfig) // Panic's if there is any error
	if len(apiForTellerConfig.RewardConfig.LadderLine) == 0 {
		apiForTellerConfig.RewardConfig.LadderLine = []int{0}
	}
	if len(apiForTellerConfig.RewardConfig.PromoterRatio) == 0 {
		apiForTellerConfig.RewardConfig.PromoterRatio = []float64{0.05}
	}
	if len(apiForTellerConfig.RewardConfig.SuperiorPromoterRatio) == 0 {
		apiForTellerConfig.RewardConfig.SuperiorPromoterRatio = []float64{0.03}
	}
	//	fmt.Printf("%+v\n", apiForTellerConfig)
	if len(apiForTellerConfig.RewardConfig.LadderLine) != len(apiForTellerConfig.RewardConfig.PromoterRatio) || len(apiForTellerConfig.RewardConfig.LadderLine) != len(apiForTellerConfig.RewardConfig.SuperiorPromoterRatio) {
		panic("RewardConfig LadderLine, PromoterRatio, SuperiorPromoterRatio length not same")
	}
	initApiForTellerConfig = true
	return apiForTellerConfig
}
