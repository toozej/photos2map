package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetEnvVars(t *testing.T) {
	tests := []struct {
		name         string
		mockEnv      map[string]string
		mockEnvFile  string
		expectError  bool
		expectDebug  bool
		expectDir    string
		expectOutput string
	}{
		{
			name: "Valid environment variables",
			mockEnv: map[string]string{
				"PHOTOS2MAP_DEBUG":  "true",
				"PHOTOS2MAP_DIR":    "/test/dir",
				"PHOTOS2MAP_OUTPUT": "gpx",
			},
			expectError:  false,
			expectDebug:  true,
			expectDir:    "/test/dir",
			expectOutput: "gpx",
		},
		{
			name:         "Valid .env file",
			mockEnvFile:  "PHOTOS2MAP_DEBUG=true\nPHOTOS2MAP_DIR=/env/dir\nPHOTOS2MAP_OUTPUT=html\n",
			expectError:  false,
			expectDebug:  true,
			expectDir:    "/env/dir",
			expectOutput: "html",
		},
		{
			name:         "No environment variables or .env file - use defaults",
			expectError:  false,
			expectDebug:  false,
			expectDir:    ".",
			expectOutput: "html",
		},
		{
			name: "Environment variable overrides .env file",
			mockEnv: map[string]string{
				"PHOTOS2MAP_DEBUG":  "false",
				"PHOTOS2MAP_DIR":    "/env/override",
				"PHOTOS2MAP_OUTPUT": "gpx",
			},
			mockEnvFile:  "PHOTOS2MAP_DEBUG=true\nPHOTOS2MAP_DIR=/file/dir\nPHOTOS2MAP_OUTPUT=html\n",
			expectError:  false,
			expectDebug:  false,
			expectDir:    "/env/override",
			expectOutput: "gpx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original directory and change to temp directory
			originalDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get current directory: %v", err)
			}

			// Save original environment variables
			originalEnvVars := map[string]string{
				"PHOTOS2MAP_DEBUG":  os.Getenv("PHOTOS2MAP_DEBUG"),
				"PHOTOS2MAP_DIR":    os.Getenv("PHOTOS2MAP_DIR"),
				"PHOTOS2MAP_OUTPUT": os.Getenv("PHOTOS2MAP_OUTPUT"),
			}
			defer func() {
				for key, value := range originalEnvVars {
					if value != "" {
						os.Setenv(key, value)
					} else {
						os.Unsetenv(key)
					}
				}
			}()

			tmpDir := t.TempDir()
			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("Failed to change to temp directory: %v", err)
			}
			defer func() {
				if err := os.Chdir(originalDir); err != nil {
					t.Errorf("Failed to restore original directory: %v", err)
				}
			}()

			// Clear environment variables first
			for key := range originalEnvVars {
				os.Unsetenv(key)
			}

			// Create .env file if applicable
			if tt.mockEnvFile != "" {
				envPath := filepath.Join(tmpDir, ".env")
				if err := os.WriteFile(envPath, []byte(tt.mockEnvFile), 0644); err != nil {
					t.Fatalf("Failed to write mock .env file: %v", err)
				}
			}

			// Set mock environment variables (these should override .env file)
			for key, value := range tt.mockEnv {
				os.Setenv(key, value)
			}

			// Call function
			conf := GetEnvVars()

			// Verify output
			if conf.Debug != tt.expectDebug {
				t.Errorf("expected debug %t, got %t", tt.expectDebug, conf.Debug)
			}
			if conf.Dir != tt.expectDir {
				t.Errorf("expected dir %q, got %q", tt.expectDir, conf.Dir)
			}
			if conf.Output != tt.expectOutput {
				t.Errorf("expected output %q, got %q", tt.expectOutput, conf.Output)
			}
		})
	}
}
