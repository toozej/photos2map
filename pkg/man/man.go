// Package man provides manual page generation functionality for the photos2map application.
//
// This package creates Unix manual pages (man pages) for the photos2map CLI
// application using the mango-cobra library. It generates properly formatted
// man pages that can be viewed with the standard `man` command on Unix systems.
//
// The package integrates with cobra commands to automatically generate documentation
// from command structures, descriptions, and usage information. Generated man pages
// follow standard Unix conventions for section 1 (user commands).
//
// Key features:
//   - Automatic man page generation from cobra command definitions
//   - Standard roff formatting for compatibility with man command
//   - Hidden command integration (not shown in help but available for internal use)
//   - Error handling for generation failures
//
// Example usage:
//
//	import "github.com/toozej/photos2map/pkg/man"
//
//	// Add man command to root command
//	rootCmd.AddCommand(man.NewManCmd())
//
//	// Generate man pages:
//	// ./photos2map man > photos2map.1
package man

import (
	"fmt"
	"os"

	mcoral "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"
)

// NewManCmd creates and returns a new cobra command for generating manual pages.
//
// This function constructs a hidden cobra command that generates Unix manual pages
// for the photos2map application. The command traverses the root command tree
// and creates comprehensive documentation in standard roff format.
//
// Command characteristics:
//   - Use: "man" - the command name for invocation
//   - Hidden: true - not shown in help output but available for use
//   - Args: cobra.NoArgs - accepts no command-line arguments
//   - SilenceUsage: true - suppresses usage on errors
//   - DisableFlagsInUseLine: true - cleaner usage line display
//
// The generated man page includes:
//   - Command descriptions and usage patterns
//   - Flag and option documentation
//   - Subcommand information
//   - Standard man page sections (NAME, SYNOPSIS, DESCRIPTION, etc.)
//
// Returns:
//   - *cobra.Command: A configured cobra command for man page generation
//
// Errors:
//   - Returns error if man page generation fails
//   - Returns error if output writing fails
//
// Example:
//
//	// Create and add man command
//	manCmd := man.NewManCmd()
//	rootCmd.AddCommand(manCmd)
//
//	// Usage from command line:
//	// ./photos2map man > photos2map.1
//	// man ./photos2map.1
func NewManCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "man",
		Short:                 "Generates photos2map's command line manpages",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Hidden:                true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			manPage, err := mcoral.NewManPage(1, cmd.Root())
			if err != nil {
				return err
			}

			_, err = fmt.Fprint(os.Stdout, manPage.Build(roff.NewDocument()))
			return err
		},
	}

	return cmd
}
