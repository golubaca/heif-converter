package main

import (
	"bytes"
	"flag"
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"time"

	"github.com/golubaca/goheif"
)

func main() {
	imagePath := flag.String("image", "", "Path to HEIF/HEIC image to convert")
	flag.Parse()
	fmt.Printf("Parsing %s\n", *imagePath)
	startTime := time.Now()
	goheif.SafeEncoding = true
	fin := imagePath
	fi, err := os.Open(*fin)
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()

	img, err := goheif.Decode(fi)
	if err != nil {
		log.Fatalf("Failed to parse %s: %v\n", *fin, err)
	}

	fout := fmt.Sprintf("%s.jpg", *imagePath)

	fo, err := os.OpenFile(fout, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to create output file %s: %v\n", fout, err)
	}
	defer fo.Close()

	err = jpeg.Encode(fo, img, nil)
	if err != nil {
		log.Fatalf("Failed to encode %s: %v\n", fout, err)
	}
	conversionTime := time.Since(startTime)

	log.Printf("Convert %s successfully in %d ms\n", *fin, conversionTime.Milliseconds())
}

func insertExif(jpegData []byte, exifData []byte) []byte {
	// Find the start of image marker
	soi := []byte{0xFF, 0xD8}

	// If the JPEG doesn't start with SOI marker, return original
	if len(jpegData) < 2 || !bytes.Equal(jpegData[0:2], soi) {
		return jpegData
	}

	// Create EXIF segment
	exifHeader := []byte{0xFF, 0xE1}
	exifSize := uint16(len(exifData) + 2) // +2 for size bytes
	exifSizeBytes := []byte{byte(exifSize >> 8), byte(exifSize & 0xFF)}

	// Combine all parts
	result := make([]byte, 0, len(jpegData)+len(exifData)+4)
	result = append(result, soi...)
	result = append(result, exifHeader...)
	result = append(result, exifSizeBytes...)
	result = append(result, exifData...)
	result = append(result, jpegData[2:]...)

	return result
}
