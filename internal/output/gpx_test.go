package output

import (
	"os"
	"testing"

	"github.com/go-echarts/go-echarts/v2/opts"
)

// TestGenerateGPX checks if a valid GPX file is generated.
func TestGenerateGPX(t *testing.T) {
	gpsData := []opts.GeoData{
		{Name: "Image1", Value: []float64{-0.1276, 51.5074}},
		{Name: "Image2", Value: []float64{2.3522, 48.8566}},
	}

	GenerateGPX(gpsData)

	// Check if the GPX file is created
	if _, err := os.Stat("out/output.gpx"); os.IsNotExist(err) {
		t.Fatalf("Expected output.gpx to be generated, but it does not exist")
	}

	// Clean up after test
	os.Remove("out/output.gpx")
}
