package service_test

import (
	"os"
	"testing"

	"github.com/juandr89/delivery-notifier-buyer/server"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.SetConfigName("config")

		configFile, err := os.Create("config.yaml")
		if err != nil {
			t.Fatalf("failed to create config file: %v", err)
		}
		defer os.Remove("config.yaml")

		_, err = configFile.WriteString(`port: "8080"
api_key: "test-api-key"`)
		if err != nil {
			t.Fatalf("failed to write to config file: %v", err)
		}
		configFile.Close()

		cfg, err := server.LoadConfig()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "8080", cfg.Port)
		assert.Equal(t, "test-api-key", cfg.APIKey)
	})

	t.Run("NotFound", func(t *testing.T) {
		viper.SetConfigName("non_existing_config")
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")

		cfg, err := server.LoadConfig()
		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "error reading config")
	})
}
