package config

import (
	"github.com/enhanced-tools/errors"
	"github.com/spf13/viper"
)

type Config struct {
	Address  string   `mapstructure:"address"`
	Postgres Postgres `mapstructure:"db"`
}

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

func validatePostgres(p *Postgres) error {
	if p.Host == "" {
		return errors.New("postgres host is required")
	}
	if p.Port == 0 {
		return errors.New("postgres port is required")
	}
	if p.User == "" {
		return errors.New("postgres user is required")
	}
	if p.Password == "" {
		return errors.New("postgres password is required")
	}
	if p.Database == "" {
		return errors.New("postgres database is required")
	}
	return nil
}

func ValidateConfig(cfg *Config) error {
	return validatePostgres(&cfg.Postgres)
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/app/")
	viper.AddConfigPath("../../")
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Enhance(err)
	}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Enhance(err)
	}
	if err := ValidateConfig(cfg); err != nil {
		return nil, errors.Enhance(err)
	}
	return cfg, nil
}
