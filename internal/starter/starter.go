package starter

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	username string
)

func Run(cmd *cobra.Command, args []string) {
	getEnvVars()

	fmt.Println("Hello from ", username)
}
