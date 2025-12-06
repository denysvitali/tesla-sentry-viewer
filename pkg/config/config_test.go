package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	// Test default config
	cfg := NewConfig()

	assert.Equal(t, "output", cfg.OutputDir)
	assert.Equal(t, "ffmpeg", cfg.FFmpegPath)
	assert.Equal(t, 4, cfg.MaxWorkers)
	assert.Equal(t, "127.0.0.1:8150", cfg.ServerAddress)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.False(t, cfg.Debug)
}

func TestConfigFromEnvironment(t *testing.T) {
	// Set environment variables with TESLA_SENTRY_ prefix
	os.Setenv("TESLA_SENTRY_DEBUG", "true")
	os.Setenv("TESLA_SENTRY_OUTPUT", "/custom/output")
	os.Setenv("TESLA_SENTRY_FFMPEG", "/usr/local/bin/ffmpeg")
	os.Setenv("TESLA_SENTRY_MAX_WORKERS", "8")
	os.Setenv("TESLA_SENTRY_SERVER_LISTEN", "0.0.0.0:8080")
	os.Setenv("TESLA_SENTRY_LOG_LEVEL", "debug")

	// Test config with environment variables
	cfg := NewConfigFromViper()

	assert.True(t, cfg.Debug)
	assert.Equal(t, "/custom/output", cfg.OutputDir)
	assert.Equal(t, "/usr/local/bin/ffmpeg", cfg.FFmpegPath)
	assert.Equal(t, 8, cfg.MaxWorkers)
	assert.Equal(t, "0.0.0.0:8080", cfg.ServerAddress)
	assert.Equal(t, "debug", cfg.LogLevel)

	// Clean up
	os.Unsetenv("TESLA_SENTRY_DEBUG")
	os.Unsetenv("TESLA_SENTRY_OUTPUT")
	os.Unsetenv("TESLA_SENTRY_FFMPEG")
	os.Unsetenv("TESLA_SENTRY_MAX_WORKERS")
	os.Unsetenv("TESLA_SENTRY_SERVER_LISTEN")
	os.Unsetenv("TESLA_SENTRY_LOG_LEVEL")
}

func TestEnsureOutputDirectory(t *testing.T) {
	// Test with existing directory
	cfg := &Config{OutputDir: "output"}
	err := EnsureOutputDirectory(cfg)
	assert.NoError(t, err)

	// Test with new directory
	cfg = &Config{OutputDir: "test_output_dir"}
	err = EnsureOutputDirectory(cfg)
	assert.NoError(t, err)

	// Clean up
	os.RemoveAll("test_output_dir")
}