package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type (
	Config struct {
		HTTP `mapstructure:"http"`
		Log  `mapstructure:"log"`
		PG   `mapstructure:"postgres"`
	}

	HTTP struct {
		Adress string
	}

	Log struct {
		Level string `mapstructure:"level"`
	}

	PG struct {
		MaxPoolSize int `mapstructure:"max_pool_size"`
		Conn        string
	}
)

func LoadConfig(configPath string) (config *Config, err error) {
	// load .env file
	err = godotenv.Load(".env")
	if err != nil {
		return &Config{}, fmt.Errorf("error loading .env file: %w", err)
	}

	// load yaml file
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return &Config{}, fmt.Errorf("error reading config file: %w", err)
	}

	if err = viper.Unmarshal(&config); err != nil {
		return &Config{}, fmt.Errorf("error unmarshaling config: %w", err)
	}

	config.Conn = os.Getenv("POSTGRES_CONN")
	config.Adress = os.Getenv("SERVER_ADDRESS")

	return config, nil
}
