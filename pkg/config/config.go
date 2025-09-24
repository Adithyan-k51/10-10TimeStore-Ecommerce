package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost           string `mapstructure:"DB_HOST" validate:"required"`
	DBName           string `mapstructure:"DB_NAME" validate:"required"`
	DBUser           string `mapstructure:"DB_USER" validate:"required"`
	DBPort           string `mapstructure:"DB_PORT" validate:"required"`
	DBPassword       string `mapstructure:"DB_PASSWORD" validate:"required"`
	AUTHTOCKEN       string `mapstructure:"TWILIO_AUTHTOCKEN" validate:"required"`
	ACCOUNTSID       string `mapstructure:"TWILIO_ACCOUNT_SID" validate:"required"`
	SERVICES_ID      string `mapstructure:"TWILIO_SERVICES_ID" validate:"required"`
	RAZOR_PAY_KEY    string `mapstructure:"RAZOR_PAY_KEY"`
	RAZOR_PAY_SECRET string `mapstructure:"RAZOR_PAY_SECRET"`
}

var envs = []string{
	"DB_HOST", "DB_NAME", "DB_USER", "DB_PORT", "DB_PASSWORD",
	"TWILIO_AUTHTOCKEN", "TWILIO_ACCOUNT_SID", "TWILIO_SERVICES_ID", //twilio
	"RAZOR_PAY_KEY", "RAZOR_PAY_SECRET", //razor
}

func LoadConfig() (Config, error) {
	var config Config

	// Set default values
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_NAME", "ecommerce")
	viper.SetDefault("DB_PASSWORD", "1234")

	// Try to load from .env file
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Warning: .env file not found or error reading it: %v\n", err)
	}

	// Check environment variables
	for _, env := range envs {
		// First check actual env variable
		if value := os.Getenv(env); value != "" {
			viper.Set(env, value)
		}

		// Bind to viper
		if err := viper.BindEnv(env); err != nil {
			return config, fmt.Errorf("failed to bind env var %s: %v", env, err)
		}
	}

	// Unmarshal config
	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// Validate config
	validate := validator.New()
	if err := validate.Struct(&config); err != nil {
		return config, fmt.Errorf("config validation failed: %v", err)
	}

	// Validate Twilio credentials format
	if len(config.ACCOUNTSID) < 34 || !strings.HasPrefix(config.ACCOUNTSID, "AC") {
		return config, fmt.Errorf("invalid TWILIO_ACCOUNT_SID format")
	}

	if len(config.AUTHTOCKEN) < 32 {
		return config, fmt.Errorf("invalid TWILIO_AUTHTOCKEN format")
	}

	if len(config.SERVICES_ID) < 34 || !strings.HasPrefix(config.SERVICES_ID, "VA") {
		return config, fmt.Errorf("invalid TWILIO_SERVICES_ID format")
	}

	return config, nil
}

func GetConfig() Config {
	config, err := LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return config
}
