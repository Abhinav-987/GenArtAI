package store

import (
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func ImageDownloader(url string) (string, error) {
	URL := url
	r, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	imageDir := "./public/images"
	fname := path.Base(url)
	filePath := filepath.Join(imageDir, fname)
	f, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, r.Body)
	if err != nil {
		return "", err
	}
	return filePath, nil
}
