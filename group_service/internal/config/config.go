package config

import (
	"github.com/num30/config"
)

type Config struct {
	ListenAddress string   `default:":8081"`
	LogLevel      string   `default:"info" flag:"loglevel"`
	DB            Database `default:"{}"`
	DebugMode     bool     `flag:"debug"`
}

type Database struct {
	Host     string `default:"localhost" validate:"required"`
	Password string `default:"group_service" validate:"required" envvar:"DB_PASS"`
	DbName   string `default:"group_service"`
	Username string `default:"group_service"`
	Port     int    `default:"5435"`
}

func MustBuild(cfgFile string) *Config {
	var conf Config
	err := config.NewConfReader(cfgFile).Read(&conf)
	if err != nil {
		panic(err)
	}

	return &conf
}
