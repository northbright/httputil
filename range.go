package httputil

import (
	"errors"
	"net/http"
)

const (
	keyAcceptRanges = "Accept-Ranges"
)

var (
	// ErrNot200or206 represents the error that the status code is not 200 or 206.
	ErrNot200or206 = errors.New("status code is not 200 or 206")
)

// IsRangeSupported returns if range is supported by the server on the URL.
func IsRangeSupported(url string) (bool, error) {
	// Do HTTP request with HEAD method(without body)
	resp, err := http.Head(url)
	if err != nil {
		return false, err
	}

	// Check if status code is 200 or 206.
	if resp.StatusCode != 200 && resp.StatusCode != 206 {
		return false, ErrNot200or206
	}

	if resp.Header.Get(keyAcceptRanges) != "bytes" {
		return false, nil
	}
	return true, nil
}
