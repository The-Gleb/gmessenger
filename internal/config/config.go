package config

import (
	"github.com/num30/config"
)

type Config struct {
	RunAddress string   `default:":8080"`
	LogLevel   string   `default:"info" flag:"loglevel"`
	DB         Database `default:"{}"`
	DebugMode  bool     `flag:"debug"`
}

type Database struct {
	Host     string `default:"localhost" validate:"required"`
	Password string `default:"gmessenger_gateway" validate:"required" envvar:"DB_PASS"`
	DbName   string `default:"gmessenger_gateway"`
	Username string `default:"gmessenger_gateway"`
	Port     int    `default:"5434"`
}

func BuildConfig(cfgFile string) *Config {
	var conf Config
	err := config.NewConfReader(cfgFile).Read(&conf)
	if err != nil {
		panic(err)
	}

	return &conf
}
