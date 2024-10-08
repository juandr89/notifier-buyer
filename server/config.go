package server

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Port                  string                `mapstructure:"port"`
	APIKey                string                `mapstructure:"api_key"`
	NotificationSender    string                `mapstructure:"notification_sender"`
	SMTPConfig            SMTPConfig            `mapstructure:"smtp"`
	RedisConfig           RedisConfig           `mapstructure:"redis"`
	ForecastServiceConfig ForecastServiceConfig `mapstructure:"forecast_service"`
}

type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type ForecastServiceConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type SendGridConfig struct {
	APIKey string `mapstructure:"api_key"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func LoadConfig() (*Config, error) {

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	//viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshal config: %w", err)
	}

	return &cfg, nil
}
