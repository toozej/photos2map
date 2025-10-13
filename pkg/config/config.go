// Package config provides secure configuration management for the photos2map application.
//
// This package handles loading configuration from environment variables and .env files
// with built-in security measures to prevent path traversal attacks. It uses the
// github.com/caarlos0/env library for environment variable parsing and
// github.com/joho/godotenv for .env file loading.
//
// The configuration loading follows a priority order:
//  1. CLI flags (highest priority)
//  2. Environment variables
//  3. .env file in current working directory
//  4. Default values (if any)
//
// Security features:
//   - Path traversal protection for .env file loading
//   - Secure file path resolution using filepath.Abs and filepath.Rel
//   - Validation against directory traversal attempts
//
// Example usage:
//
//	import "github.com/toozej/photos2map/pkg/config"
//
//	func main() {
//		conf := config.GetEnvVars()
//		fmt.Printf("Debug: %t\n", conf.Debug)
//	}
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Config represents the application configuration structure.
//
// This struct defines all configurable parameters for the photos2map
// application. Fields are tagged with struct tags that correspond to
// environment variable names for automatic parsing.
//
// Currently supported configuration:
//   - Debug: Enable debug logging, loaded from PHOTOS2MAP_DEBUG env var
//   - Dir: Directory to scan for images, loaded from PHOTOS2MAP_DIR env var
//   - Output: Output format (html or gpx), loaded from PHOTOS2MAP_OUTPUT env var
//
// Example:
//
//	type Config struct {
//		Debug  bool   `env:"PHOTOS2MAP_DEBUG"`
//		Dir    string `env:"PHOTOS2MAP_DIR" envDefault:"."`
//		Output string `env:"PHOTOS2MAP_OUTPUT" envDefault:"html"`
//	}
type Config struct {
	// Debug specifies whether to enable debug-level logging.
	// It is loaded from the PHOTOS2MAP_DEBUG environment variable.
	// If not set, defaults to false.
	Debug bool `env:"PHOTOS2MAP_DEBUG"`

	// Dir specifies the directory to scan for images.
	// It is loaded from the PHOTOS2MAP_DIR environment variable.
	// If not set, defaults to current directory (".").
	Dir string `env:"PHOTOS2MAP_DIR" envDefault:"."`

	// Output specifies the output format (html or gpx).
	// It is loaded from the PHOTOS2MAP_OUTPUT environment variable.
	// If not set, defaults to "html".
	Output string `env:"PHOTOS2MAP_OUTPUT" envDefault:"html"`
}

// GetEnvVars loads and returns the application configuration from environment
// variables and .env files with comprehensive security validation.
//
// This function performs the following operations:
//  1. Securely determines the current working directory
//  2. Constructs and validates the .env file path to prevent traversal attacks
//  3. Loads .env file if it exists in the current directory
//  4. Parses environment variables into the Config struct
//  5. Returns the populated configuration
//
// Security measures implemented:
//   - Path traversal detection and prevention using filepath.Rel
//   - Absolute path resolution for secure path operations
//   - Validation against ".." sequences in relative paths
//   - Safe file existence checking before loading
//
// The function will terminate the program with os.Exit(1) if any critical
// errors occur during configuration loading, such as:
//   - Current directory access failures
//   - Path traversal attempts detected
//   - .env file parsing errors
//   - Environment variable parsing failures
//
// Returns:
//   - Config: A populated configuration struct with values from environment
//     variables and/or .env file
//
// Example:
//
//	// Load configuration
//	conf := config.GetEnvVars()
//
//	// Use configuration
//	if conf.Debug {
//		fmt.Println("Debug mode enabled")
//	}
func GetEnvVars() Config {
	// Get current working directory for secure file operations
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %s\n", err)
		os.Exit(1)
	}

	// Construct secure path for .env file within current directory
	envPath := filepath.Join(cwd, ".env")

	// Ensure the path is within our expected directory (prevent traversal)
	cleanEnvPath, err := filepath.Abs(envPath)
	if err != nil {
		fmt.Printf("Error resolving .env file path: %s\n", err)
		os.Exit(1)
	}
	cleanCwd, err := filepath.Abs(cwd)
	if err != nil {
		fmt.Printf("Error resolving current directory: %s\n", err)
		os.Exit(1)
	}
	relPath, err := filepath.Rel(cleanCwd, cleanEnvPath)
	if err != nil || strings.Contains(relPath, "..") {
		fmt.Printf("Error: .env file path traversal detected\n")
		os.Exit(1)
	}

	// Load .env file if it exists
	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			fmt.Printf("Error loading .env file: %s\n", err)
			os.Exit(1)
		}
	}

	// Parse environment variables into config struct
	var conf Config
	if err := env.Parse(&conf); err != nil {
		fmt.Printf("Error parsing environment variables: %s\n", err)
		os.Exit(1)
	}

	return conf
}
