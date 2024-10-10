package output

import (
	"os"
	"testing"

	"github.com/go-echarts/go-echarts/v2/opts"
)

// TestGenerateMap checks that the HTML map file is created correctly.
func TestGenerateMap(t *testing.T) {
	gpsData := []opts.GeoData{
		{Name: "Image1", Value: []float64{-0.1276, 51.5074}},
		{Name: "Image2", Value: []float64{2.3522, 48.8566}},
	}

	GenerateMap(gpsData)

	// Check if the map file is created
	if _, err := os.Stat("out/map.html"); os.IsNotExist(err) {
		t.Fatalf("Expected map.html to be generated, but it does not exist")
	}

	// Clean up after test
	os.Remove("out/map.html")
}
