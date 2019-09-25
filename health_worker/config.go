package health_worker

import (
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

const AuthTypeHTTP = "http"
const AuthTypeToken = "token"

func BindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		switch v.Kind() {
		case reflect.Struct:
			BindEnvs(v.Interface(), append(parts, tv)...)
		default:
			viper.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
}

func NewConfig(confPath string) (Config, error) {
	var conf Config
	defaultConfig(&conf)

	viper.SetConfigFile(confPath)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("PIR5")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return conf, err
	}

	BindEnvs(conf)

	if err := viper.Unmarshal(&conf); err != nil {
		return conf, err
	}

	return conf, nil
}

type Config struct {
	WorkerID     int      `mapstructure:"worker_id"`
	PollInterval int      `mapstructure:"poll_interval"`
	Concurrency  int      `mapstructure:"concurrency"`
	Listen       string   `mapstructure:"listen"`
	DB           database `mapstructure:"database"`
	Redis        redis    `mapstructure:"redis"`
	PdnsAPI      pdnsAPI  `mapstructure:"pdnsAPI"`
	Auth         auth    `mapstructure:"auth"`
}

type database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DBName   string `mapstructure:"dbname"`
	UserName string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	Password string `mapstructure:"password"`
	TTL      int    `mapstructure:"ttl"`
	PoolSize int    `mapstructure:"pool_size"`
}

type pdnsAPI struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
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
	c.PdnsAPI.Host = "127.0.0.1"
	c.PdnsAPI.Port = 8080
}

type auth struct {
	AuthType string   `mapstructure:"auth_type"` // token or http
	Tokens   []string `mapstructure:"tokens"`
	HttpAuth httpAuth `mapstructure:"http_auth"`
}

func (c Config) IsTokenAuth() bool {
	return c.Auth.AuthType == AuthTypeToken
}

func (c Config) IsHTTPAuth() bool {
	return c.Auth.AuthType == AuthTypeHTTP
}
