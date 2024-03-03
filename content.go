package httputil

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

const (
	keyAcceptRanges  = "Accept-Ranges"
	keyContentLength = "Content-Length"
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

// getResp returns the response, length of the content and if range is supported by the server.
func getResp(uri string, method string) (*http.Response, uint64, bool, error) {
	if method != "HEAD" && method != "GET" {
		return nil, 0, false, ErrMethodNotHeadOrGet
	}

	// Create an HTTP client.
	client := http.Client{}
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return nil, 0, false, err
	}

	// Do HTTP request.
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, false, err
	}

	// Check if status code is 200.
	if resp.StatusCode != 200 {
		return nil, 0, false, ErrNot200
	}

	supported := false
	if resp.Header.Get(keyAcceptRanges) == "bytes" {
		supported = true
	}

	l, _ := strconv.ParseUint(resp.Header.Get(keyContentLength), 10, 64)

	return resp, l, supported, nil
}

// Len returns the length of the content and if range is supported by the server.
func Len(uri string) (l uint64, isRangeSupported bool, err error) {
	resp, l, isRangeSupported, err := getResp(uri, "HEAD")
	if err != nil {
		return l, false, err
	}
	defer resp.Body.Close()

	return l, isRangeSupported, nil
}

// GetResp returns the response, length of content and if range is supported by the server.
func GetResp(uri string) (resp *http.Response, l uint64, isRangeSupported bool, err error) {
	return getResp(uri, "GET")
}

// SetRangeHeader adds the range key-value pair to the header.
// If isEndIgnored is true, the syntax is "bytes=start-".
func SetRangeHeader(header http.Header, start, end uint64, isEndIgnored bool) {
	bytesRange := ""
	if !isEndIgnored {
		bytesRange = fmt.Sprintf("bytes=%d-%d", start, end)
	} else {
		bytesRange = fmt.Sprintf("bytes=%d-", start)
	}
	header.Add("range", bytesRange)
}

// getRespOfRange returns the response and the size of the partial content.
// If isEndIgnored is true, the range header uses "bytes=start-" syntax.
func getRespOfRange(uri string, method string, start, end uint64, isEndIgnored bool) (*http.Response, uint64, error) {
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
	SetRangeHeader(req.Header, start, end, isEndIgnored)

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
	l, _ := strconv.ParseUint(str, 10, 64)

	return resp, l, err
}

// LenOfRange returns the size of the partial content.
// If isEndIgnored is true, the range header uses "bytes=start-" syntax.
func LenOfRange(uri string, start, end uint64, isEndIgnored bool) (uint64, error) {
	resp, l, err := getRespOfRange(uri, "HEAD", start, end, isEndIgnored)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return l, nil
}

// LenOfRangeStart returns the size of the partial content.
// The range header uses "bytes=start-" syntax.
func LenOfRangeStart(uri string, start uint64) (uint64, error) {
	return LenOfRange(uri, start, 0, true)
}

// GetRespOfRange returns the response and the size of the partial content.
// If isEndIgnored is true, the range header uses "bytes=start-" syntax.
func GetRespOfRange(uri string, start, end uint64, isEndIgnored bool) (*http.Response, uint64, error) {
	return getRespOfRange(uri, "GET", start, end, isEndIgnored)
}

// GetRespOfRangeStart returns the response and the size of the partial content.
// If isEndIgnored is true, the range header uses "bytes=start-" syntax.
func GetRespOfRangeStart(uri string, start uint64) (*http.Response, uint64, error) {
	return GetRespOfRange(uri, start, 0, true)
}
