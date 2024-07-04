package store

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	cld "github.com/Abhinav-987/GenArtAI/pkg/cloudinary"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

func ImageDownloader(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %w", err)
	}
	publicID := uuid.New().String()
	imageReader := bytes.NewReader(imageData)
	uploadResult, err := cld.Cld.Upload.Upload(context.Background(), imageReader, uploader.UploadParams{
		PublicID: publicID,
		Folder:   "abhinav",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to Cloudinary: %w", err)
	}

	return uploadResult.SecureURL, nil
}
