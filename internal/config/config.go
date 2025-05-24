package config

import (
	"github.com/enhanced-tools/errors"
	"github.com/spf13/viper"
)

type Config struct {
	Address  string   `mapstructure:"address"`
	Postgres Postgres `mapstructure:"db"`
	MCP      MCP      `mapstructure:"mcp"`
}

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type MCP struct {
	Address string `mapstructure:"address"`
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

// LoadConfig loads configuration, optionally from a specific file path.
func LoadConfigFromFile(configPath string) (*Config, error) {
	cfg := &Config{
		Address: ":8080", // Default address
		Postgres: Postgres{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			Database: "postgres",
		},
		MCP: MCP{
			Address: "stdio",
		},
	}

	viper.SetConfigType("yaml")
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/app/")
		viper.AddConfigPath("../../")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, errors.Enhance(err)
		}
		// If config file is not found, we'll use the defaults
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Enhance(err)
	}

	if err := ValidateConfig(cfg); err != nil {
		return nil, errors.Enhance(err)
	}
	return cfg, nil
}

// LoadConfig is preserved for backward compatibility (no argument)
func LoadConfig() (*Config, error) {
	return LoadConfigFromFile("")
}
