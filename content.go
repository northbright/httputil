package httputil

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

var (
	// status code is not 200
	ErrNot200 = errors.New("status code is not 200")

	// content range is not set correctly
	ErrIncorrectRange = errors.New("incorrect range")

	// status code is not 200, 416 or 206
	ErrNot200or416or206 = errors.New("status code is not 200, 416 or 206")

	// range header is not supported by the server
	ErrRangeNotSupported = errors.New("range header is not supported by the server")

	// request method is not HEAD or GET
	ErrMethodNotHeadOrGet = errors.New("request method is not HEAD or GET")
)

// getResp returns the HTTP response, the size of the content and if range header is supported by the server.
// -1 size indicates that the size is unknown.
func getResp(uri string, method string) (resp *http.Response, size int64, rangeIsSupported bool, err error) {
	if method != "HEAD" && method != "GET" {
		return nil, -1, false, ErrMethodNotHeadOrGet
	}

	// Create an HTTP client.
	client := http.Client{}
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return nil, -1, false, err
	}

	// Do HTTP request.
	if resp, err = client.Do(req); err != nil {
		return nil, -1, false, err
	}

	// Check if status code is 200.
	if resp.StatusCode != 200 {
		return nil, -1, false, ErrNot200
	}

	size = int64(-1)
	str := resp.Header.Get("Content-Length")
	if str != "" {
		if size, err = strconv.ParseInt(str, 10, 64); err != nil {
			resp.Body.Close()
			return nil, -1, false, err
		}
	}

	// Check if range header is supported.
	rangeIsSupported = false
	if resp.Header.Get("Accept-Ranges") == "bytes" {
		rangeIsSupported = true
	}

	return resp, size, rangeIsSupported, nil
}

// Size returns the size of the content and if range header is supported by the server.
// -1 size indicates the size is unknown.
func Size(uri string) (size int64, rangeIsSupported bool, err error) {
	resp, size, rangeIsSupported, err := getResp(uri, "HEAD")
	if err != nil {
		return -1, false, err
	}
	defer resp.Body.Close()

	return size, rangeIsSupported, nil
}

// GetResp returns the HTTP response, the size of the content and if range header is supported by the server.
// -1 size indicates that the size is unknown.
func GetResp(uri string) (resp *http.Response, size int64, rangeIsSupported bool, err error) {
	return getResp(uri, "GET")
}

// SetRangeHeader adds the range key-value pair to the header.
// If endIsIgnored is true, the syntax is "bytes=start-".
func SetRangeHeader(header http.Header, start, end int64, endIsIgnored bool) {
	bytesRange := ""
	if !endIsIgnored {
		bytesRange = fmt.Sprintf("bytes=%d-%d", start, end)
	} else {
		bytesRange = fmt.Sprintf("bytes=%d-", start)
	}
	header.Add("range", bytesRange)
}

// getRespOfRange returns the response and the size of the partial content.
// If endIsIgnored is true, the range header uses "bytes=start-" syntax.
func getRespOfRange(uri string, method string, start, end int64, endIsIgnored bool) (*http.Response, int64, error) {
	if method != "HEAD" && method != "GET" {
		return nil, 0, ErrMethodNotHeadOrGet
	}

	// Create an HTTP client.
	client := http.Client{}
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return nil, 0, err
	}

	// Set the range header to resume the downloading if need.
	SetRangeHeader(req.Header, start, end, endIsIgnored)

	// Do HTTP request.
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	// Check the status code.
	if resp.StatusCode != 206 {
		switch resp.StatusCode {
		case 200:
			// Server may return 200 if range is not supported.
			return nil, 0, ErrRangeNotSupported
		case 416:
			return nil, 0, ErrIncorrectRange
		default:
			return nil, 0, ErrNot200or416or206
		}
	}

	// Get the remote file size.
	str := resp.Header.Get("Content-Length")
	size, _ := strconv.ParseInt(str, 10, 64)

	return resp, size, err
}

// SizeOfRange returns the size of the partial content.
// If endIsIgnored is true, the range header uses "bytes=start-" syntax.
func SizeOfRange(uri string, start, end int64, endIsIgnored bool) (int64, error) {
	resp, l, err := getRespOfRange(uri, "HEAD", start, end, endIsIgnored)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return l, nil
}

// SizeOfRangeStart returns the size of the partial content.
// The range header uses "bytes=start-" syntax.
func SizeOfRangeStart(uri string, start int64) (int64, error) {
	return SizeOfRange(uri, start, 0, true)
}

// GetRespOfRange returns the response and the size of the partial content.
// If endIsIgnored is true, the range header uses "bytes=start-" syntax.
func GetRespOfRange(uri string, start, end int64, endIsIgnored bool) (*http.Response, int64, error) {
	return getRespOfRange(uri, "GET", start, end, endIsIgnored)
}

// GetRespOfRangeStart returns the response and the size of the partial content.
func GetRespOfRangeStart(uri string, start int64) (*http.Response, int64, error) {
	return GetRespOfRange(uri, start, 0, true)
}
