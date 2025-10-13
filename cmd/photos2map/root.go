// Package cmd provides the command-line interface for the photos2map application.
//
// This package implements the main CLI functionality using the cobra library,
// handling command parsing, flag management, and orchestrating the core photo
// processing workflow. It serves as the entry point for user interactions with
// the photos2map application.
//
// The package provides:
//   - Root command definition and configuration
//   - CLI flag and argument parsing
//   - Integration with configuration management (env/godotenv)
//   - Coordination between EXIF processing, GPS extraction, and output generation
//   - Debug logging configuration
//   - Subcommand integration (version, man pages)
//
// Command structure:
//   - photos2map: Main command for processing photos and generating maps/GPX
//   - photos2map version: Display version and build information
//   - photos2map man: Generate manual pages
//
// Configuration priority (highest to lowest):
//  1. CLI flags and arguments
//  2. Environment variables (PHOTOS2MAP_*)
//  3. .env file in current directory
//  4. Default values
//
// Example usage:
//
//	import "github.com/toozej/photos2map/cmd/photos2map"
//
//	// Execute the CLI
//	cmd.Execute()
//
// Command line examples:
//
//	# Process photos in current directory, output HTML map
//	photos2map
//
//	# Process specific directory, output GPX file
//	photos2map --dir ./vacation-photos --output gpx
//
//	# Enable debug logging
//	photos2map --debug --dir ./photos
package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/toozej/photos2map/internal/extract"
	"github.com/toozej/photos2map/internal/output"
	"github.com/toozej/photos2map/pkg/config"
	"github.com/toozej/photos2map/pkg/man"
	"github.com/toozej/photos2map/pkg/version"
)

// debugFlag is a global variable to track debug mode state.
//
// This variable is used to store the debug flag value since we no longer use
// viper for flag binding. It's set during command initialization and used
// in the PersistentPreRun function to configure logging levels.
var debugFlag bool

// rootCmd defines the base command when called without any subcommands.
//
// This is the main entry point for the photos2map CLI application. It processes
// photos in a specified directory, extracts GPS coordinates from EXIF data,
// and generates either an HTML map or GPX file based on user preferences.
//
// The command supports the following workflow:
//  1. Load configuration from environment variables and .env files
//  2. Parse CLI flags and arguments (with priority over config)
//  3. Configure debug logging if requested
//  4. Extract GPS data from photo EXIF information
//  5. Generate output in the requested format (HTML map or GPX file)
//
// Command characteristics:
//   - Use: "photos2map" - the main command name
//   - Args: cobra.ExactArgs(0) - accepts no positional arguments
//   - Flags: --dir (input directory), --output (format), --debug (logging)
//   - PersistentPreRun: Configures logging and loads configuration
//   - Run: Executes the main photo processing workflow
var rootCmd = &cobra.Command{
	Use:              "photos2map",
	Short:            "Generate a GPX file from photos EXIF data",
	Long:             `Generates a map on a HTML page or GPX file from GPS coordinates in images`,
	Args:             cobra.ExactArgs(0),
	PersistentPreRun: rootCmdPreRun,
	Run:              run,
}

// rootCmdPreRun configures the application environment before command execution.
//
// This function is called before the main command runs and handles:
//   - Loading configuration from environment variables and .env files
//   - Setting up debug logging based on CLI flags or configuration
//   - Ensuring proper logging level configuration
//
// The function follows configuration priority order:
//  1. CLI debug flag (highest priority)
//  2. PHOTOS2MAP_DEBUG environment variable
//  3. Debug setting from .env file
//  4. Default (no debug logging)
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command line arguments (unused in this function)
func rootCmdPreRun(cmd *cobra.Command, args []string) {
	// Load configuration from environment variables and .env file
	conf := config.GetEnvVars()

	// Set debug level based on CLI flag (highest priority) or config
	if debugFlag || conf.Debug {
		log.SetLevel(log.DebugLevel)
	}
}

// Execute runs the root command and handles any execution errors.
//
// This is the main entry point for the CLI application, typically called from
// main.go. It executes the root cobra command and handles any errors that occur
// during command parsing or execution.
//
// The function will:
//   - Parse command line arguments and flags
//   - Execute the appropriate command (root, version, man, etc.)
//   - Print error messages and exit with code 1 if execution fails
//   - Exit normally (code 0) if execution succeeds
//
// Error handling:
//   - Prints error message to stdout if command execution fails
//   - Calls os.Exit(1) to terminate with error status
//   - Does not return if an error occurs
//
// Example:
//
//	func main() {
//		cmd.Execute() // This will handle all CLI interaction
//	}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// init initializes the command structure and configures flags and subcommands.
//
// This function is called automatically when the package is imported and sets up:
//   - CLI flags and their bindings to global variables
//   - Subcommand registration (version, man pages)
//   - Command structure and relationships
//
// Flag configuration:
//   - --debug/-d: Boolean flag for debug logging (bound to debugFlag variable)
//   - --dir/-i: String flag for input directory (default: current directory)
//   - --output/-o: String flag for output format (default: "html")
//
// Subcommands added:
//   - version: Display version and build information
//   - man: Generate manual pages for the application
//
// The function uses cobra's flag binding mechanisms to connect CLI flags
// with the command execution logic.
func init() {
	// create rootCmd-level flags
	rootCmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "d", false, "Enable debug-level logging")
	rootCmd.Flags().StringP("dir", "i", ".", "Directory to scan for images")
	rootCmd.Flags().StringP("output", "o", "html", "Output format: html or gpx")

	// add sub-commands
	rootCmd.AddCommand(
		man.NewManCmd(),
		version.Command(),
	)
}

// run executes the core photo processing workflow.
//
// This function implements the main application logic:
//  1. Loads configuration from environment variables and .env files
//  2. Retrieves CLI flag values with proper priority handling
//  3. Extracts GPS coordinates from photo EXIF data
//  4. Generates output in the requested format (HTML map or GPX file)
//
// Configuration priority (highest to lowest):
//   - CLI flags (if explicitly set by user)
//   - Environment variables (PHOTOS2MAP_DIR, PHOTOS2MAP_OUTPUT)
//   - .env file values
//   - Default values
//
// The function processes all images in the specified directory, extracts GPS
// coordinates from EXIF data, and generates either an interactive HTML map
// or a GPX file suitable for GPS devices and mapping applications.
//
// Parameters:
//   - cmd: The cobra command being executed (used to access flag values)
//   - args: Command line arguments (unused in current implementation)
//
// Output:
//   - HTML map file (default): Interactive map showing photo locations
//   - GPX file: GPS exchange format file for GPS devices
//   - Console message if no GPS data found in images
//
// Example workflow:
//  1. User runs: photos2map --dir ./vacation --output gpx
//  2. Function loads config, processes ./vacation directory
//  3. Extracts GPS data from JPEG EXIF information
//  4. Generates vacation.gpx file with GPS waypoints
func run(cmd *cobra.Command, args []string) {
	// Load configuration from environment variables and .env file
	conf := config.GetEnvVars()

	// Get CLI flag values (these take priority over config)
	dir, _ := cmd.Flags().GetString("dir")
	outputType, _ := cmd.Flags().GetString("output")

	// Use config values if CLI flags are at default values
	if !cmd.Flags().Changed("dir") && conf.Dir != "." {
		dir = conf.Dir
	}
	if !cmd.Flags().Changed("output") && conf.Output != "html" {
		outputType = conf.Output
	}

	gpsData := extract.ExtractGPSData(dir)

	if len(gpsData) > 0 {
		if outputType == "gpx" {
			output.GenerateGPX(gpsData)
		} else {
			output.GenerateMap(gpsData)
		}
	} else {
		fmt.Println("No GPS data found in the images.")
	}
}
