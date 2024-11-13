package main

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/golubaca/goheif"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

type FileConversionInfo struct {
	OriginalFileName string
	OriginalFileSize int64
	NewFileName      string
	NewFileSize      int64
	ConversionTime   time.Duration
	Thumbnail        []byte
}

type ConversionResult struct {
	Error     error
	TotalTime time.Duration
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ConvertMultipleFiles(filePaths []string) ConversionResult {
	startTime := time.Now()
	resultsChan := make(chan FileConversionInfo, len(filePaths))
	var wg sync.WaitGroup

	for i, path := range filePaths {
		wg.Add(1)
		go func(index int, filePath string) {
			defer wg.Done()
			result, err := a.ConvertFile(filePath)
			if err != nil {
				runtime.EventsEmit(a.ctx, "conversion_error", map[string]interface{}{
					"path":  filePath,
					"error": err.Error(),
				})
				return
			}
			resultsChan <- result
			runtime.EventsEmit(a.ctx, "conversion_progress", result)
		}(i, path)
	}

	// Close channel after all goroutines complete
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results from channel
	var results []FileConversionInfo
	for result := range resultsChan {
		results = append(results, result)
	}

	conversionTime := time.Since(startTime)
	finalResult := ConversionResult{
		Error:     nil,
		TotalTime: time.Duration(conversionTime.Milliseconds()),
	}

	runtime.EventsEmit(a.ctx, "conversion_complete", finalResult)
	return finalResult
}

func (a *App) OpenFilePicker() ([]string, error) {
	filters := []runtime.FileFilter{
		{DisplayName: "HEIF Images", Pattern: "*.heif;*.heic"},
	}
	paths, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		Filters: filters,
	})
	if err != nil {
		return nil, err
	}

	// Ensure full paths are returned
	fullPaths := make([]string, len(paths))
	for i, path := range paths {
		fullPaths[i], err = filepath.Abs(path)
		if err != nil {
			return nil, err
		}
	}

	return fullPaths, nil
}

func (a *App) ConvertFile(filePath string) (FileConversionInfo, error) {
	startTime := time.Now()
	fi, err := os.Open(filePath)
	if err != nil {
		return FileConversionInfo{}, err
	}
	defer fi.Close()

	exif, err := goheif.ExtractExif(fi)
	if err != nil {
		return FileConversionInfo{}, err
	}

	img, err := goheif.Decode(fi)
	if err != nil {
		return FileConversionInfo{}, err
	}

	// Create JPEG in memory
	var jpegBuffer bytes.Buffer
	err = jpeg.Encode(&jpegBuffer, img, nil)
	if err != nil {
		return FileConversionInfo{}, err
	}

	// Get JPEG bytes and insert EXIF
	jpegBytes := jpegBuffer.Bytes()
	// Replace with actual EXIF data
	finalJpegBytes := insertExif(jpegBytes, exif)

	// Create output directory
	convertDir := filepath.Join(filepath.Dir(filePath), "convert")
	err = os.MkdirAll(convertDir, os.ModePerm)
	if err != nil {
		return FileConversionInfo{}, err
	}

	// Write final JPEG with EXIF to file
	outputFilePath := filepath.Join(convertDir, filepath.Base(filePath)+".jpg")
	err = os.WriteFile(outputFilePath, finalJpegBytes, 0644)
	if err != nil {
		return FileConversionInfo{}, err
	}

	originalFileInfo, err := os.Stat(filePath)
	if err != nil {
		return FileConversionInfo{}, err
	}

	newFileInfo, err := os.Stat(outputFilePath)
	if err != nil {
		return FileConversionInfo{}, err
	}

	conversionTime := time.Since(startTime)

	return FileConversionInfo{
		OriginalFileName: filePath,
		OriginalFileSize: originalFileInfo.Size(),
		NewFileName:      outputFilePath,
		NewFileSize:      newFileInfo.Size(),
		ConversionTime:   time.Duration(conversionTime.Milliseconds()),
		Thumbnail:        nil,
	}, nil
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

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
