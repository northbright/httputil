package httputil

import (
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
)

// GetFileNameFromURL returns the downloadable file name from target URL.
func GetFileNameFromURL(targetURL string) (string, error) {
	if targetURL == "" {
		return "", fmt.Errorf("empty target URL")
	}

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return "", err
	}

	fileName := filepath.Base(parsedURL.Path)
	if fileName == "." {
		return "", fmt.Errorf("parsedURL.Path err: %v", parsedURL.Path)
	}

	return fileName, nil
}

// GetFileNameFromResponse returns downloadable file name in the HTTP response.
func GetFileNameFromResponse(resp *http.Response) (string, error) {
	contentDisposition := resp.Header.Get("Content-Disposition")

	if contentDisposition == "" {
		return "", fmt.Errorf("Content-Disposition is empty")
	}

	_, params, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		return "", err
	}

	fileName, ok := params["filename"]
	if !ok {
		return "", fmt.Errorf("no filename param in Content-Disposition")
	}

	return fileName, nil
}

// GetFileName returns the downloadable file name.
// It'll try to get the name in the HTTP response firstly.
// If it fails, then it'll try to parse the URL to get the file name.
func GetFileName(targetURL string) (string, error) {
	// Try to do HTTP request(HEAD).
	resp, err := http.Head(targetURL)
	if err != nil {
		// Try to get the file name in URL.
		return GetFileNameFromURL(targetURL)
	}
	defer resp.Body.Close()

	// Try to get the file name in the response.
	fileName, err := GetFileNameFromResponse(resp)
	if err != nil {
		// Try to get the file name in the URL.
		return GetFileNameFromURL(targetURL)
	}

	return fileName, nil
}
