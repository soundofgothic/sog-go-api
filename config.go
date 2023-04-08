package main

import (
	"github.com/enhanced-tools/errors"
	"github.com/spf13/viper"
)

type Config struct {
	Address  string   `yaml:"address"`
	Postgres Postgres `yaml:"postgres"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
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

func setDefaults() {
	viper.SetDefault("address", ":3000")
	viper.SetDefault("postgres.host", "localhost")
	viper.SetDefault("postgres.port", 5432)
	viper.SetDefault("postgres.user", "postgres")
	viper.SetDefault("postgres.password", "postgres")
	viper.SetDefault("postgres.database", "postgres")
}

func ValidateConfig(cfg *Config) error {
	return validatePostgres(&cfg.Postgres)
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	setDefaults()
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
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
