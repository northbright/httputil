package httputil

import (
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
)

// GetFileNameFromURL gets the downloadable file name from target URL.
//
//   Params:
//     targetURL: Target URL.
//   Return:
//     fileName: File name in the URL.
func GetFileNameFromURL(targetURL string) (fileName string, err error) {
	var parsedURL *url.URL

	if targetURL == "" {
		return "", fmt.Errorf("Empty target URL.")
	}

	if parsedURL, err = url.Parse(targetURL); err != nil {
		return "", err
	}
	if fileName = filepath.Base(parsedURL.Path); fileName == "." {
		return "", fmt.Errorf("parsedURL.Path err: %v\n", parsedURL.Path)
	}

	return fileName, nil
}

// GetFileNameFromResponse detects downladable file name in the HTTP response.
//
//   Params:
//     resp: http.Response returned from http.Do(), Head(), Get()...
//   Return:
//     fileName: File name.
func GetFileNameFromResponse(resp *http.Response) (fileName string, err error) {
	var ok bool
	contentDisposition := resp.Header.Get("Content-Disposition")

	if contentDisposition == "" {
		return "", fmt.Errorf("Content-Disposition is empty.")
	}

	_, params, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		return "", err
	}

	if fileName, ok = params["filename"]; !ok {
		return "", fmt.Errorf("No filename param in Content-Disposition.")
	}

	return fileName, nil

}

// GetFileName detects / gets the downladable file name in the given target URL.
//
//   Params:
//     targetURL: Target URL.
//   Return:
//     fileName: File name.
func GetFileName(targetURL string) (fileName string, err error) {
	var resp *http.Response

	if resp, err = http.Head(targetURL); err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Step 1. Try to detect the downloadable file name from HTTP response.
	if fileName, err = GetFileNameFromResponse(resp); err != nil {
		// Step 2. Try to get the downloadable file name in the URL.
		return GetFileNameFromURL(targetURL)
	}

	return fileName, nil
}
