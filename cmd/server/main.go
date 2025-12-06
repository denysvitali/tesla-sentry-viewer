package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/denysvitali/tesla-sentry-viewer/pkg/config"
	"github.com/denysvitali/tesla-sentry-viewer/pkg/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const SoftwareName = "tesla-sentry-viewer"

var Version = "dev"

var cfgFile string

var rootCmd = &cobra.Command{
	Use:     "tesla-sentry-viewer <directory>",
	Short:   "Tesla Sentry clip viewer and server",
	Long:    `A web server that serves Tesla Sentry Mode video clips for viewing and processing.`,
	Version: Version,
	Args:    cobra.ExactArgs(1),
	Run:     runServer,
}

func init() {
	cobra.OnInitialize(initConfig)

	// Persistent flags (available to all commands)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tesla-sentry-viewer.yaml)")
	rootCmd.PersistentFlags().BoolP("debug", "D", false, "Enable debug logging")
	rootCmd.PersistentFlags().StringP("listen", "l", "127.0.0.1:8150", "Server listen address")
	rootCmd.PersistentFlags().StringP("output", "o", "output", "Output directory for processed files")
	rootCmd.PersistentFlags().String("ffmpeg", "ffmpeg", "Path to ffmpeg binary")

	// Bind flags to viper
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("server.listen", rootCmd.PersistentFlags().Lookup("listen"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("ffmpeg", rootCmd.PersistentFlags().Lookup("ffmpeg"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Warning: could not find home directory:", err)
		} else {
			viper.AddConfigPath(home)
		}
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".tesla-sentry-viewer")
	}

	// Environment variables
	viper.SetEnvPrefix("TESLA_SENTRY")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Read config file (ignore if not found)
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func runServer(cmd *cobra.Command, args []string) {
	directory := args[0]

	// Build config from viper
	cfg := &config.Config{
		Debug:         viper.GetBool("debug"),
		OutputDir:     viper.GetString("output"),
		FFmpegPath:    viper.GetString("ffmpeg"),
		MaxWorkers:    viper.GetInt("server.max_workers"),
		ServerAddress: viper.GetString("server.listen"),
		LogLevel:      viper.GetString("log_level"),
	}

	// Apply defaults for values not set
	if cfg.MaxWorkers == 0 {
		cfg.MaxWorkers = 4
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	// Setup logger
	logger := logrus.New()
	if cfg.Debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	// Validate directory
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		logger.Fatalf("Directory does not exist: %s", directory)
	}

	// Ensure output directory exists
	if err := config.EnsureOutputDirectory(cfg); err != nil {
		logger.Fatalf("Failed to setup output directory: %v", err)
	}

	logger.Infof("Starting %s %s", SoftwareName, Version)
	logger.Infof("Serving directory: %s", directory)
	logger.Infof("Listening on: %s", cfg.ServerAddress)
	logger.Infof("Output directory: %s", cfg.OutputDir)

	// Create and start server
	s, err := server.NewWithConfig(directory, cfg)
	if err != nil {
		logger.Fatalf("unable to create server: %v", err)
	}

	s.SetLogger(logger)
	err = s.Listen(cfg.ServerAddress)
	if err != nil {
		logger.Fatalf("unable to start server: %v", err)
	}
}

func main() {
	// Set version for --version flag
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate(fmt.Sprintf("%s {{.Version}}\n", SoftwareName))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
