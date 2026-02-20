package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

func ProcessImage(inputPath string) (string, error) {
	src, err := imaging.Open(inputPath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %w", err)
	}

	maxWidth := 800
	quality := 80
	uploadDir := os.Getenv("UPLOAD_DIR")

	if uploadDir == "" {
		uploadDir = "./public/temp"
	}
	if w := os.Getenv("IMG_MAX_WIDTH"); w != "" {
		if parsed, err := strconv.Atoi(w); err == nil {
			maxWidth = parsed
		}
	}
	if q := os.Getenv("IMG_QUALITY"); q != "" {
		if parsed, err := strconv.Atoi(q); err == nil {
			quality = parsed
		}
	}

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	resized := imaging.Resize(src, maxWidth, 0, imaging.Lanczos)

	outputFilename := fmt.Sprintf("recipe_%s.jpg", uuid.New().String()[:8])
	outputPath := filepath.Join(uploadDir, outputFilename)

	err = imaging.Save(resized, outputPath, imaging.JPEGQuality(quality))
	if err != nil {
		return "", fmt.Errorf("failed to save processed image: %w", err)
	}

	_ = os.Remove(inputPath)

	return outputFilename, nil
}
