package exif

import (
	"os"
	"testing"
)

// Test for successful EXIF extraction
// TODO mock or create valid test image file with EXIF GPS data
// func TestExtractEXIF_Success(t *testing.T) {
// 	// Create a temporary file to simulate a valid image file
// 	file, err := os.CreateTemp("", "valid.jpg")
// 	if err != nil {
// 		t.Fatalf("failed to create temp file: %v", err)
// 	}
// 	defer os.Remove(file.Name()) // Clean up temp file after test
//
// 	lat, lon, err := ExtractEXIF(file.Name())
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
//
// 	if lat != 37.7749 {
// 		t.Errorf("expected latitude 37.7749, got %v", lat)
// 	}
// 	if lon != -122.4194 {
// 		t.Errorf("expected longitude -122.4194, got %v", lon)
// 	}
// }

// Test for file open failure
func TestExtractEXIF_FileOpenError(t *testing.T) {
	// Try opening a non-existent file
	_, _, err := ExtractEXIF("nonexistent.jpg")
	if err == nil {
		t.Error("expected an error when opening non-existent file, got none")
	}
}

// Test for EXIF decode failure
func TestExtractEXIF_DecodeError(t *testing.T) {
	// Create a temporary file to simulate an invalid image file
	file, err := os.CreateTemp("", "invalid.jpg")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name()) // Clean up temp file after test

	_, _, err = ExtractEXIF(file.Name())
	if err == nil {
		t.Error("expected an error during EXIF decoding, got none")
	}
}
