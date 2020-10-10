package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"time"
)

// App config struct
type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	MongoDB  MongoDB
	Cookie   Cookie
	Store    Store
	Session  Session
	Metrics  Metrics
}

// Server config struct
type ServerConfig struct {
	Port              string
	PprofPort         string
	Mode              string
	JwtSecretKey      string
	CookieName        string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	SSL               bool
	CtxDefaultTimeout time.Duration
}

// Postgresql config struct
type PostgresConfig struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PostgresqlSSLMode  bool
	PgDriver           string
}

// Redis config struct
type RedisConfig struct {
	RedisAddr      string
	RedisPassword  string
	RedisDb        string
	RedisDefaultdb string
}

// MongoDB config struct
type MongoDB struct {
	MongoURI string
}

// Cookie config struct
type Cookie struct {
	Name     string
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

type Session struct {
	Prefix string
	Name   string
	Expire int
}

type Metrics struct {
	Url         string
	ServiceName string
}

type Store struct {
	ImagesFolder string
}

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		} else {
			return nil, err
		}
	}

	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
