package config

import (
	"fmt"

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
		return fmt.Errorf("postgres host is required")
	}
	if p.Port == 0 {
		return fmt.Errorf("postgres port is required")
	}
	if p.User == "" {
		return fmt.Errorf("postgres user is required")
	}
	if p.Password == "" {
		return fmt.Errorf("postgres password is required")
	}
	if p.Database == "" {
		return fmt.Errorf("postgres database is required")
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
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
		// If config file is not found, we'll use the defaults
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := ValidateConfig(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	return cfg, nil
}

// LoadConfig is preserved for backward compatibility (no argument)
func LoadConfig() (*Config, error) {
	return LoadConfigFromFile("")
}
