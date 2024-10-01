package starter

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

// mockExit allows us to override os.Exit in tests.
func mockExit(code int) {
	panic("os.Exit called")
}

func TestGetEnvVars_Success(t *testing.T) {
	// Setup environment variables for a successful run
	os.Setenv("USERNAME", "testuser")

	// Clear Viper config in case it is retaining state from previous tests
	viper.Reset()

	// Defer unsetting environment variables to avoid pollution
	defer os.Unsetenv("USERNAME")

	// Override the exit function to prevent the test from exiting
	exit = mockExit
	defer func() { exit = os.Exit }() // Reset after the test

	// Call the function and check that it runs without exiting
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected no panic, got %v", r)
		}
	}()

	// Run the function
	getEnvVars()

	// Validate that the expected values were retrieved
	if username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", username)
	}
}

func TestGetEnvVars_MissingUsername(t *testing.T) {
	// Setup environment variables with missing username
	os.Unsetenv("USERNAME")

	// Clear Viper config in case it is retaining state from previous tests
	viper.Reset()

	// Defer unsetting environment variables to avoid pollution
	defer os.Unsetenv("USERNAME")

	// Override the exit function to prevent the test from exiting
	exit = mockExit
	defer func() { exit = os.Exit }() // Reset after the test

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic due to missing USERNAME")
		}
	}()

	// Run the function (expecting it to panic due to missing username)
	getEnvVars()
}
