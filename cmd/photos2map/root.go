package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/toozej/photos2map/internal/extract"
	"github.com/toozej/photos2map/internal/output"
	"github.com/toozej/photos2map/pkg/man"
	"github.com/toozej/photos2map/pkg/version"
)

var rootCmd = &cobra.Command{
	Use:              "photos2map",
	Short:            "Generate a GPX file from photos EXIF data",
	Long:             `Generates a map on a HTML page or GPX file from GPS coordinates in images`,
	Args:             cobra.ExactArgs(0),
	PersistentPreRun: rootCmdPreRun,
	Run:              run,
}

func rootCmdPreRun(cmd *cobra.Command, args []string) {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return
	}
	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func init() {
	_, err := maxprocs.Set()
	if err != nil {
		log.Error("Error setting maxprocs: ", err)
	}

	// create rootCmd-level flags
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug-level logging")
	rootCmd.Flags().StringP("dir", "i", ".", "Directory to scan for images")
	rootCmd.Flags().StringP("output", "o", "html", "Output format: html or gpx")
	_ = viper.BindPFlag("dir", rootCmd.Flags().Lookup("dir"))
	_ = viper.BindPFlag("output", rootCmd.Flags().Lookup("output"))

	// add sub-commands
	rootCmd.AddCommand(
		man.NewManCmd(),
		version.Command(),
	)
}

// Core functionality to process the images and output either an HTML map or GPX file
func run(cmd *cobra.Command, args []string) {
	dir := viper.GetString("dir")
	outputType := viper.GetString("output")
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
