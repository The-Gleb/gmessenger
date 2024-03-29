package config

import (
	"time"

	"github.com/num30/config"
)

type Config struct {
	RunAddress      string        `default:":8080" envvar:"GATEWAY_GATEWAY_ADDR"`
	LogLevel        string        `default:"info" flag:"loglevel" envvar:"LOGLEVEL"`
	TokenTTL        time.Duration `default:"24h"`
	GroupServerHost string        `default:"localhost" envvar:"GROUP_SERVICE_HOST"`
	GroupServerPort int           `default:"8081" envvar:"GROUP_SERVICE_PORT"`
	DB              Database      `default:"{}"`
	DebugMode       bool          `flag:"debug"`
}

type Database struct {
	Host     string `default:"localhost" validate:"required" envvar:"GATEWAY_DB_HOST"`
	Port     int    `default:"5434" envvar:"GATEWAY_DB_PORT"`
	Password string `default:"gateway_db" validate:"required" envvar:"GATEWAY_DB_PASS"`
	DbName   string `default:"gateway_db" envvar:"GATEWAY_DB_NAME"`
	Username string `default:"gateway_db" envvar:"GATEWAY_DB_USERNAME"`
}

func MustBuild(cfgFile string) *Config {
	var conf Config
	err := config.NewConfReader(cfgFile).Read(&conf)
	if err != nil {
		panic(err)
	}

	return &conf
}
