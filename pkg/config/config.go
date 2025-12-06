package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Debug          bool
	OutputDir      string
	FFmpegPath     string
	MaxWorkers     int
	ServerAddress  string
	LogLevel       string
}

func NewConfig() *Config {
	return &Config{
		Debug:          getViperBool("debug", false),
		OutputDir:      getViperString("output", "output"),
		FFmpegPath:     getViperString("ffmpeg", "ffmpeg"),
		MaxWorkers:     getViperInt("max_workers", 4),
		ServerAddress:  getViperString("server.listen", "127.0.0.1:8150"),
		LogLevel:       getViperString("log_level", "info"),
	}
}

func NewConfigFromViper() *Config {
	// Initialize Viper with defaults
	viper.SetDefault("debug", false)
	viper.SetDefault("output", "output")
	viper.SetDefault("ffmpeg", "ffmpeg")
	viper.SetDefault("server.listen", "127.0.0.1:8150")
	viper.SetDefault("max_workers", 4)
	viper.SetDefault("log_level", "info")

	// Read from environment variables
	viper.SetEnvPrefix("TESLA_SENTRY")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	return &Config{
		Debug:          viper.GetBool("debug"),
		OutputDir:      viper.GetString("output"),
		FFmpegPath:     viper.GetString("ffmpeg"),
		MaxWorkers:     viper.GetInt("max_workers"),
		ServerAddress:  viper.GetString("server.listen"),
		LogLevel:       viper.GetString("log_level"),
	}
}

func getViperString(key, defaultValue string) string {
	if viper.IsSet(key) {
		return viper.GetString(key)
	}
	return defaultValue
}

func getViperBool(key string, defaultValue bool) bool {
	if viper.IsSet(key) {
		return viper.GetBool(key)
	}
	return defaultValue
}

func getViperInt(key string, defaultValue int) int {
	if viper.IsSet(key) {
		return viper.GetInt(key)
	}
	return defaultValue
}

// Backward compatibility functions
func getStringEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return strings.ToLower(value) == "true" || value == "1"
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}

func EnsureOutputDirectory(config *Config) error {
	if config.OutputDir == "" {
		config.OutputDir = "output"
	}

	// Create output directory if it doesn't exist
	if _, err := os.Stat(config.OutputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
			return err
		}
	}

	// Make it absolute path
	absPath, err := filepath.Abs(config.OutputDir)
	if err != nil {
		return err
	}
	config.OutputDir = absPath

	return nil
}