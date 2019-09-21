package health_worker

import (
	"github.com/BurntSushi/toml"
)

func NewConfig(confPath string) (Config, error) {
	var conf Config
	defaultConfig(&conf)

	if _, err := toml.DecodeFile(confPath, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}

type Config struct {
	WorkerID     int
	PollInterval int
	Concurrency  int
	Listen       string   `toml:"listen"`
	DB           database `toml:"database"`
	Redis        redis    `toml:"redis"`
	TokenAuth    *tokenAuth
}

type database struct {
	Host     string
	Port     int
	DBName   string `toml:"dbname"`
	UserName string `toml:"username"`
	Password string
}

type redis struct {
	Host     string
	Port     int
	DB       int
	Password string
	TTL      int
	PoolSize int
}

func defaultConfig(c *Config) {
	c.Listen = "0.0.0.0:8080"
	c.PollInterval = 10
	c.Concurrency = 10000
	c.DB.Host = "localhost"
	c.DB.Port = 3306
	c.DB.UserName = "root"
	c.DB.DBName = "health_worker"
	c.Redis.Host = "localhost"
	c.Redis.Port = 6379
	c.Redis.PoolSize = 10
	c.Redis.DB = 0
}

type tokenAuth struct {
	Tokens []string
}

func (c Config) IsTokenAuth() bool {
	return c.TokenAuth != nil
}
