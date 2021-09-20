package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Environment int

const (
	DevelopmentEnv Environment = iota
	ProductionEnv
)

type (
	AppConfig struct {
		Name string      `toml:"name"`
		Env  Environment `toml:"env"`
		Port string      `toml:"port"`
	}

	DatabaseConfig struct {
		Driver    string `toml:"driver"`
		Name      string `toml:"name"`
		Host      string `toml:"host"`
		Username  string `toml:"username"`
		Password  string `toml:"password"`
		ParseTime bool   `toml:"parse_time"`
	}

	LoggerConfig struct {
		MaxAge     int    `toml:"max_age"`
		MaxSize    int    `toml:"max_size"`
		MaxBackups int    `toml:"max_backups"`
		Directory  string `toml:"directory"`
	}

	CacheConfig struct {
		Driver string `toml:"driver"`
	}

	RedisConfig struct {
		Addresses []string `toml:"addresses"`
		Password  string   `toml:"password"`
		Database  int      `toml:"database"`
		Timeout   int64    `toml:"timeout"` // in seconds
	}

	Config struct {
		App      AppConfig
		Database DatabaseConfig
		Logger   LoggerConfig
		Redis    RedisConfig
		Cache    CacheConfig
	}
)

func Load(configFile string) (*Config, error) {
	var cfg *Config

	if _, err := toml.DecodeFile(configFile, &cfg); err != nil {
		return nil, fmt.Errorf("error decoding config file: %s", err.Error())
	}

	return cfg, nil
}
