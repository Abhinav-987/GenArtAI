package store

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Abhinav-987/GenArtAI/pkg/cloudinary"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

func ImageDownloader(url string) (string, error) {
	// Download the image from the given URL
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %w", err)
	}
	tmpFile, err := os.CreateTemp("", "image-*.jpg")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.Write(imageData); err != nil {
		return "", fmt.Errorf("failed to write image data to temp file: %w", err)
	}
	uploadResult, err := cloudinary.Cld.Upload.Upload(context.Background(), tmpFile.Name(), uploader.UploadParams{})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to Cloudinary: %w", err)
	}
	return uploadResult.SecureURL, nil
}
