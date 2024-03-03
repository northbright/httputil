package httputil

import (
	"errors"
	"net/http"
	"strconv"
)

const (
	keyAcceptRanges  = "Accept-Ranges"
	keyContentLength = "Content-Length"
)

var (
	// ErrNot200or206 represents the error that the status code is not 200 or 206.
	ErrNot200or206 = errors.New("status code is not 200 or 206")
)

// Len returns the length of content and if range is supported by the server on the URL.
func Len(uri string) (uint64, bool, error) {
	// Do HTTP request with HEAD method(without body)
	resp, err := http.Head(uri)
	if err != nil {
		return 0, false, err
	}
	defer resp.Body.Close()

	// Check if status code is 200 or 206.
	if resp.StatusCode != 200 && resp.StatusCode != 206 {
		return 0, false, ErrNot200or206
	}

	supported := false
	if resp.Header.Get(keyAcceptRanges) == "bytes" {
		supported = true
	}

	l, _ := strconv.ParseUint(resp.Header.Get(keyContentLength), 10, 64)

	return l, supported, nil
}
