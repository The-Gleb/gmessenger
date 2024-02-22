package config

import (
	"log"

	"github.com/num30/config"
)

type Config struct {
	ListenAddress string   `default:":8081" envvar:"GROUP_LISTEN_ADDR"`
	LogLevel      string   `default:"info" flag:"loglevel"`
	DB            Database `default:"{}"`
	DebugMode     bool     `flag:"debug"`
}

type Database struct {
	Host     string `default:"localhost" validate:"required" envvar:"GROUP_DB_HOST"`
	Password string `default:"group_db" validate:"required" envvar:"GROUP_DB_PASS"`
	DbName   string `default:"group_db" envvar:"GROUP_DB_NAME"`
	Username string `default:"group_db" envvar:"GROUP_DB_USERNAME"`
	Port     int    `default:"5433" envvar:"GROUP_DB_PORT"`
}

func MustBuild(cfgFile string) *Config {
	var conf Config
	err := config.NewConfReader(cfgFile).Read(&conf)
	if err != nil {
		log.Fatal(err)
	}

	return &conf
}
