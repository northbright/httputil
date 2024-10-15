package httputil_test

import (
	"fmt"
	"log"

	"github.com/northbright/httputil"
)

func ExampleSize() {
	uri := "https://golang.google.cn/dl/go1.19.3.darwin-arm64.pkg"

	size, rangeIsSupported, err := httputil.Size(uri)
	if err != nil {
		log.Printf("Size() error: %v", err)
		return
	}

	fmt.Printf("size: %d\nrange is supported: %v",
		size,
		rangeIsSupported)

	// Output:
	// size: 145565374
	// range is supported: true
}

func ExampleGetResp() {
	uri := "https://golang.google.cn/dl/go1.19.3.darwin-arm64.pkg"

	resp, size, rangeIsSupported, err := httputil.GetResp(uri)
	if err != nil {
		log.Printf("GetResp() error: %v", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("size: %d\nrange is supported: %v",
		size,
		rangeIsSupported)

	// Output:
	// size: 145565374
	// range is supported: true
}

func ExampleSizeOfRange() {
	uri := "https://golang.google.cn/dl/go1.19.3.darwin-arm64.pkg"

	// Get the size of the range of the content using "bytes=start-end" syntax.
	l, err := httputil.SizeOfRange(uri, 0, 99999999, false)
	if err != nil {
		log.Printf("httputil.SizeOfRange() error: %v", err)
		return
	}
	fmt.Printf("size of range: 0 - 99999999: %d\n", l)

	// Get the size of the range of the content using "bytes=start-" syntax.
	l, err = httputil.SizeOfRange(uri, 100000000, 0, true)
	if err != nil {
		log.Printf("httputil.SizeOfRange() error: %v", err)
		return
	}

	fmt.Printf("size of range: 10000000-: %d\n", l)

	// Output:
	// size of range: 0 - 99999999: 100000000
	// size of range: 10000000-: 45565374
}

func ExampleSizeOfRangeStart() {
	uri := "https://golang.google.cn/dl/go1.19.3.darwin-arm64.pkg"

	// Get the size of the range of the content using "bytes=start-" syntax.
	l, err := httputil.SizeOfRangeStart(uri, 100000000)
	if err != nil {
		log.Printf("httputil.SizeOfRangeStart() error: %v", err)
		return
	}

	fmt.Printf("size of range: 10000000-: %d\n", l)

	// Output:
	// size of range: 10000000-: 45565374
}

func ExampleGetRespOfRange() {
	uri := "https://golang.google.cn/dl/go1.19.3.darwin-arm64.pkg"

	// Get the response, size of the range of the content using "bytes=start-end" syntax.
	resp, l, err := httputil.GetRespOfRange(uri, 0, 99999999, false)
	if err != nil {
		log.Printf("httputil.GetRespOfRange() error: %v", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("size of range: 0 - 99999999: %d\n", l)

	// Get response and len of range.
	resp2, l, err := httputil.GetRespOfRange(uri, 100000000, 0, true)
	if err != nil {
		log.Printf("httputil.GetRespOfRange() error: %v", err)
		return
	}
	defer resp2.Body.Close()

	fmt.Printf("size of range: 10000000-: %d\n", l)

	// Output:
	// size of range: 0 - 99999999: 100000000
	// size of range: 10000000-: 45565374
}

func ExampleGetRespOfRangeStart() {
	uri := "https://golang.google.cn/dl/go1.19.3.darwin-arm64.pkg"

	// Get the response, size of the range of the content using "bytes=start-" syntax.
	resp, l, err := httputil.GetRespOfRangeStart(uri, 100000000)
	if err != nil {
		log.Printf("httputil.SizeOfRange() error: %v", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("size of range: 10000000-: %d\n", l)

	// Output:
	// size of range: 10000000-: 45565374
}
