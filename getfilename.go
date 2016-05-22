package httputil

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
)

// GetFileNameFromURL() gets the downloadable file name from target URL.
//
//   Params:
//     targetUrl: Target URL.
//   Return:
//     fileName: File name in the URL.
func GetFileNameFromURL(targetUrl string) (fileName string, err error) {
	if targetUrl == "" {
		return "", errors.New("Empty target URL.")
	}

	if parsedURL, err := url.Parse(targetUrl); err != nil {
		return "", err
	} else {
		if fileName = filepath.Base(parsedURL.Path); fileName == "." {
			return "", errors.New(fmt.Sprintf("parsedURL.Path err: %v\n", parsedURL.Path))
		}
	}

	return fileName, nil
}

// GetFileNameFromResponse() detects downladable file name in the HTTP response.
//
//   Params:
//     resp: http.Response returned from http.Do(), Head(), Get()...
//   Return:
//     fileName: File name.
func GetFileNameFromResponse(resp *http.Response) (fileName string, err error) {
	var ok bool
	contentDisposition := resp.Header.Get("Content-Disposition")

	if contentDisposition == "" {
		return "", errors.New("Content-Disposition is empty.")
	} else {
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err != nil {
			return "", err
		}

		if fileName, ok = params["filename"]; !ok {
			return "", errors.New("No filename param in Content-Disposition.")
		}

		return fileName, nil
	}
}

// GetFileName() detects / gets the downladable file name in the given target URL.
//
//   Params:
//     targetUrl: Target URL.
//   Return:
//     fileName: File name.
func GetFileName(targetUrl string) (fileName string, err error) {
	var resp *http.Response

	if resp, err = http.Head(targetUrl); err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Step 1. Try to detect the downloadable file name from HTTP response.
	if fileName, err = GetFileNameFromResponse(resp); err != nil {
		// Step 2. Try to get the downloadable file name in the URL.
		return GetFileNameFromURL(targetUrl)
	}

	return fileName, nil
}
