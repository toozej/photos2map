package exif

import (
	"os"

	"github.com/rwcarlsen/goexif/exif"
)

func ExtractEXIF(path string) (float64, float64, error) {
	file, err := os.Open(path) //#nosec G304
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	x, err := exif.Decode(file)
	if err != nil {
		return 0, 0, err
	}

	lat, lon, err := x.LatLong()
	return lat, lon, err
}
