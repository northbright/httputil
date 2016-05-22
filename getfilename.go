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

// GetFileName() detects / gets the downladable file name in the HTTP response.
//
//   Params:
//     targetUrl: Target URL.
//     resp: http.Response returned from http.Do(), Head(), Get()...
//   Return:
//     fileName: File name.
func GetFileName(targetUrl string, resp *http.Response) (fileName string, err error) {
	var ok bool
	contentDisposition := resp.Header.Get("Content-Disposition")

	if contentDisposition == "" {
		return GetFileNameFromURL(targetUrl)
	} else {
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err != nil {
			return "", err
		}

		if fileName, ok = params["filename"]; !ok {
			return GetFileNameFromURL(targetUrl)
		}

		return fileName, nil
	}
}
