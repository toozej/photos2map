package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/toozej/golang-starter/internal/starter"
	"github.com/toozej/golang-starter/pkg/man"
	"github.com/toozej/golang-starter/pkg/version"
)

var rootCmd = &cobra.Command{
	Use:              "golang-starter",
	Short:            "golang starter examples",
	Long:             `Golang starter template using cobra and viper modules`,
	Args:             cobra.ExactArgs(0),
	PersistentPreRun: rootCmdPreRun,
	Run:              starter.Run,
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

	// add sub-commands
	rootCmd.AddCommand(
		man.NewManCmd(),
		version.Command(),
	)
}
