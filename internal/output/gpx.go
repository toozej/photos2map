package output

import (
	"encoding/xml"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/twpayne/go-gpx"
)

// GenerateGPX creates a GPX file from the extracted GPS data.
// It takes a slice of GeoData and outputs a GPX file named `output.gpx`.
func GenerateGPX(gpsData []opts.GeoData) {
	g := gpx.GPX{
		Version: "1.1",
		Creator: "photos2map",
		Wpt:     make([]*gpx.WptType, len(gpsData)),
	}

	for i, data := range gpsData {
		// Type assert data.Value as []float64
		coords, ok := data.Value.([]float64)
		if !ok || len(coords) != 2 {
			log.Printf("Invalid GPS data for %s, skipping...", data.Name)
			continue
		}

		lat, lon := coords[1], coords[0]
		g.Wpt[i] = &gpx.WptType{
			Lat:  lat,
			Lon:  lon,
			Name: data.Name,
		}
	}

	// create output.gpx file
	file, err := os.Create("out/output.gpx")
	if err != nil {
		log.Fatalf("Error creating GPX file: %v", err)
	}
	defer file.Close()

	// Marshal the GPX struct into indented XML
	gpxData, err := xml.MarshalIndent(g, "", "  ")
	if err != nil {
		log.Errorf("Error marshalling GPX struct to XML: %v", err)
	}

	// Add the XML header and append the marshaled GPX data
	header := []byte(xml.Header)
	gpxData = append(header, gpxData...)

	// write out the gpx file
	_, err = file.Write(gpxData)
	if err != nil {
		log.Errorf("Error writing GPS data to GPX file: %v", err)
	}

	log.Println("GPX file generated successfully.")
}
