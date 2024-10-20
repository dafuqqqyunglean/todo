package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Postgres PostgresConfig
	Redis    RedisConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

func NewConfig() (Config, error) {
	err := initConfig()
	if err != nil {
		return Config{}, err
	}

	return Config{
		Postgres: PostgresConfig{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetString("db.port"),
			Username: viper.GetString("db.username"),
			DBName:   viper.GetString("db.dbname"),
			SSLMode:  viper.GetString("db.sslmode"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
		},
		Redis: RedisConfig{
			Address:  fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.GetString("redis.port")),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
		},
	}, nil
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
