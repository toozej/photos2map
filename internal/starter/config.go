package starter

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var exit = os.Exit // Use a variable for os.Exit to allow overriding in tests

// Get environment variables
func getEnvVars() {
	if _, err := os.Stat(".env"); err == nil {
		// Initialize Viper from .env file
		viper.SetConfigFile(".env") // Specify the name of your .env file

		// Read the .env file
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Error reading .env file: %s\n", err)
			exit(1)
		}
	}

	// Enable reading environment variables
	viper.AutomaticEnv()

	// get username from Viper
	username = viper.GetString("USERNAME")
	if username == "" {
		fmt.Println("username must be provided")
		exit(1)
	}
}
