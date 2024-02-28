package health_worker

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

func bindEnvs(iface interface{}, path ...string) {
	var refType reflect.Type
	var refVal reflect.Value

	if reflect.ValueOf(iface).Kind() == reflect.Ptr {
		refVal = reflect.ValueOf(iface).Elem()
		refType = refVal.Type()
	} else {
		refVal = reflect.ValueOf(iface)
		refType = refVal.Type()
	}

	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		tag, ok := field.Tag.Lookup("mapstructure")
		if !ok {
			continue // No mapstructure tag, skip this field
		}

		fieldPath := append(path, tag)
		envVar := strings.ToUpper(strings.Join(fieldPath, "_"))
		key := strings.Join(fieldPath, ".")

		// Check if the field is a struct or a pointer to a struct
		var fieldType reflect.Type
		if field.Type.Kind() == reflect.Ptr {
			// The field is a pointer, get the type it points to
			fieldType = field.Type.Elem()
		} else {
			fieldType = field.Type
		}

		for _, ak := range viper.AllKeys() {
			if ak == key {
				continue
			}
		}
		if field.Type.Kind() != reflect.Struct &&
			fieldType.Kind() != reflect.Struct &&
			field.Type.Kind() != reflect.Slice &&
			fieldType.Kind() != reflect.Map {
			viper.BindEnv(key, envVar)
		}

		// If the field value is a nil pointer, initialize it with a new zero value of the type it points to.
		if field.Type.Kind() == reflect.Ptr && refVal.Field(i).IsNil() {
			refVal.Field(i).Set(reflect.New(fieldType))
		}

		// Recursively apply environment variable bindings for nested structs
		if fieldType.Kind() == reflect.Struct {
			bindEnvs(refVal.Field(i).Interface(), fieldPath...)
		} else if field.Type.Kind() == reflect.Ptr && fieldType.Kind() == reflect.Struct {
			bindEnvs(refVal.Field(i).Elem().Interface(), fieldPath...)
		}
	}
}

func NewConfig(confPath string) (*Config, error) {
	var conf Config
	defaultConfig(&conf)

	viper.SetConfigFile(confPath)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("PIR5")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("readin config error %s", err)
	}

	bindEnvs(conf)

	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("unmarshal config error %s", err)
	}

	return &conf, nil
}

type Config struct {
	WorkerID     int      `mapstructure:"worker_id"`
	PollInterval int      `mapstructure:"poll_interval"`
	Concurrency  int      `mapstructure:"concurrency"`
	Listen       string   `mapstructure:"listen"`
	Endpoint     string   `mapstructure:"endpoint"`
	DB           database `mapstructure:"database"`
	Redis        redis    `mapstructure:"redis"`
	PdnsAPI      pdnsAPI  `mapstructure:"pdnsAPI"`
	Auth         auth     `mapstructure:"auth"`
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
	Tokens []string `mapstructure:"tokens"`
}

func (c Config) IsTokenAuth() bool {
	return len(c.Auth.Tokens) > 0
}
