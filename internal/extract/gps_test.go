package extract

import (
	"path/filepath"
	"testing"
)

// TestExtractGPSData ensures GPS coordinates are correctly extracted from test images.
func TestExtractGPSData(t *testing.T) {
	// Create a test directory with test images
	// images from: https://github.com/ianare/exif-samples
	testDir := filepath.Join("..", "testdata")

	// Call the function
	gpsData := ExtractGPSData(testDir)

	// Assert GPS data is non-empty for valid test images
	if len(gpsData) == 0 {
		t.Fatalf("Expected GPS data, but got none")
	}

	// Example: Check if first entry contains valid GPS data
	//nolint:gocritic:sloppyTypeAssert
	if gpsData[0].Name == "" || len(gpsData[0].Value.([]float64)) != 2 {
		t.Errorf("Invalid GPS data for first image: %+v", gpsData[0])
	}
}
