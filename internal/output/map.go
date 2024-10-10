// Package output provides functions to generate output formats like HTML maps or GPX files.
package output

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

// GenerateMap creates an HTML file with a world map and pins based on GPS coordinates extracted from images.
// The map is saved to "map.html".
func GenerateMap(gpsData []opts.GeoData) {
	geo := charts.NewGeo()
	geo.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "photos2map: GPS Image Map"}),
		charts.WithGeoComponentOpts(opts.GeoComponent{
			// map comes from https://github.com/echarts-maps/echarts-countries-js/tree/master/echarts-countries-js
			Map:       "USA",
			ItemStyle: &opts.ItemStyle{Color: "#006666"},
		}),
	)

	geo.AddSeries("geo", types.ChartEffectScatter, gpsData,
		charts.WithRippleEffectOpts(opts.RippleEffect{
			Period:    4,
			Scale:     6,
			BrushType: "stroke",
		}),
	)

	file, err := os.Create("out/map.html")
	if err != nil {
		log.Fatalf("Error creating map file: %v", err)
	}
	defer file.Close()

	err = geo.Render(file)
	if err != nil {
		log.Errorf("Error rendering map file to html: %v", err)
	}
	log.Println("HTML map generated successfully.")
}
