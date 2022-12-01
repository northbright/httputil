package httputil

import (
	"net/http"
)

const (
	keyAcceptRanges = "Accept-Ranges"
)

// IsRangeSupported returns if range is supported by the server on the URL.
func IsRangeSupported(url string) (bool, error) {
	resp, err := http.Head(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.Header.Get(keyAcceptRanges) != "bytes" {
		return false, nil
	}
	return true, nil
}
