package cli

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Flags represents command line flags
type Flags struct {
	ConfigFile string
	LogLevel   string
	LogFormat  string
	LogOutput  string
	ServerHost string
	ServerPort int
	Env        string
	Help       bool
	Version    bool
}

// ParseFlags parses command line flags
func ParseFlags() *Flags {
	flags := &Flags{}

	// Define flags
	pflag.StringVarP(&flags.ConfigFile, "config", "c", "", "Path to configuration file")
	pflag.StringVarP(&flags.LogLevel, "log-level", "l", "", "Log level (debug, info, warn, error)")
	pflag.StringVarP(&flags.LogFormat, "log-format", "f", "", "Log format (json, text)")
	pflag.StringVarP(&flags.LogOutput, "log-output", "o", "", "Log output (stdout, stderr, file)")
	pflag.StringVarP(&flags.ServerHost, "host", "H", "", "Server host")
	pflag.IntVarP(&flags.ServerPort, "port", "p", 0, "Server port")
	pflag.StringVarP(&flags.Env, "env", "e", "", "Environment (development, staging, production)")
	pflag.BoolVarP(&flags.Help, "help", "h", false, "Show help message")
	pflag.BoolVarP(&flags.Version, "version", "v", false, "Show version information")

	// Parse flags
	pflag.Parse()

	// Show help if requested
	if flags.Help {
		showHelp()
		os.Exit(0)
	}

	// Show version if requested
	if flags.Version {
		showVersion()
		os.Exit(0)
	}

	// Bind flags to viper
	bindFlags(flags)

	return flags
}

// bindFlags binds command line flags to viper
func bindFlags(flags *Flags) {
	if flags.ConfigFile != "" {
		viper.SetConfigFile(flags.ConfigFile)
	}

	if flags.LogLevel != "" {
		viper.Set("logging.level", flags.LogLevel)
	}

	if flags.LogFormat != "" {
		viper.Set("logging.format", flags.LogFormat)
	}

	if flags.LogOutput != "" {
		viper.Set("logging.output", flags.LogOutput)
	}

	if flags.ServerHost != "" {
		viper.Set("server.host", flags.ServerHost)
	}

	if flags.ServerPort != 0 {
		viper.Set("server.port", flags.ServerPort)
	}

	if flags.Env != "" {
		viper.Set("app.env", flags.Env)
	}
}

// showHelp displays help message
func showHelp() {
	fmt.Printf("Usage: %s [OPTIONS]\n\n", os.Args[0])
	fmt.Println("Options:")
	pflag.PrintDefaults()
	fmt.Println("\nEnvironment Variables:")
	fmt.Println("  APP_CONFIG_FILE     Path to configuration file")
	fmt.Println("  APP_LOG_LEVEL       Log level (debug, info, warn, error)")
	fmt.Println("  APP_LOG_FORMAT      Log format (json, text)")
	fmt.Println("  APP_LOG_OUTPUT      Log output (stdout, stderr, file)")
	fmt.Println("  APP_SERVER_HOST     Server host")
	fmt.Println("  APP_SERVER_PORT     Server port")
	fmt.Println("  APP_ENV             Environment (development, staging, production)")
	fmt.Println("\nConfiguration:")
	fmt.Println("  Configuration can be provided via:")
	fmt.Println("  - Command line flags (highest priority)")
	fmt.Println("  - Environment variables")
	fmt.Println("  - Configuration file")
	fmt.Println("  - Default values (lowest priority)")
}

// showVersion displays version information
func showVersion() {
	fmt.Println("YouApp v1.0.0")
	fmt.Println("A Go application with modern architecture")
}

// GetConfigFile returns the config file path from flags or environment
func GetConfigFile() string {
	if configFile := viper.GetString("config"); configFile != "" {
		return configFile
	}
	return ""
}

// GetLogLevel returns the log level from flags or environment
func GetLogLevel() string {
	return viper.GetString("logging.level")
}

// GetLogFormat returns the log format from flags or environment
func GetLogFormat() string {
	return viper.GetString("logging.format")
}

// GetLogOutput returns the log output from flags or environment
func GetLogOutput() string {
	return viper.GetString("logging.output")
}

// GetServerHost returns the server host from flags or environment
func GetServerHost() string {
	return viper.GetString("server.host")
}

// GetServerPort returns the server port from flags or environment
func GetServerPort() int {
	return viper.GetInt("server.port")
}

// GetEnv returns the environment from flags or environment
func GetEnv() string {
	return viper.GetString("app.env")
}
