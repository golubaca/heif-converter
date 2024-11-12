package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/h2non/bimg"
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

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
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

	buffer, err := os.ReadFile(filePath)
	if err != nil {
		return FileConversionInfo{}, err
	}

	image := bimg.NewImage(buffer)
	jpeg, err := image.Convert(bimg.JPEG)
	if err != nil {
		return FileConversionInfo{}, err
	}

	convertDir := filepath.Join(filepath.Dir(filePath), "convert")
	err = os.MkdirAll(convertDir, os.ModePerm)
	if err != nil {
		return FileConversionInfo{}, err
	}

	outputFilePath := filepath.Join(convertDir, filepath.Base(filePath)+".jpg")
	err = os.WriteFile(outputFilePath, jpeg, 0644)
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

	// Generate thumbnail
	thumbnail, err := bimg.NewImage(jpeg).Resize(200, 200)
	if err != nil {
		return FileConversionInfo{}, err
	}

	conversionTime := time.Since(startTime)

	return FileConversionInfo{
		OriginalFileName: filePath,
		OriginalFileSize: originalFileInfo.Size(),
		NewFileName:      outputFilePath,
		NewFileSize:      newFileInfo.Size(),
		ConversionTime:   conversionTime,
		Thumbnail:        thumbnail,
	}, nil
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
