package httputil_test

import (
	"fmt"
	"log"

	"github.com/northbright/httputil"
)

func ExampleLen() {
	uri := "https://golang.google.cn/dl/go1.19.3.darwin-arm64.pkg"

	l, supported, err := httputil.Len(uri)
	if err != nil {
		log.Printf("Len() error: %v", err)
		return
	}

	fmt.Printf("%v, %v", l, supported)

	// Output:
	// 145565374, true
}

func ExampleGetResp() {
	uri := "https://golang.google.cn/dl/go1.19.3.darwin-arm64.pkg"

	resp, l, supported, err := httputil.GetResp(uri)
	if err != nil {
		log.Printf("GetResp() error: %v", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("%v, %v", l, supported)

	// Output:
	// 145565374, true
}

func ExampleLenOfRange() {
	uri := "https://golang.google.cn/dl/go1.19.3.darwin-arm64.pkg"

	l, err := httputil.LenOfRange(uri, 0, 99999999, false)
	if err != nil {
		log.Printf("httputil.LenOfRange() error: %v", err)
		return
	}
	fmt.Printf("size of range: 0 - 99999999: %d\n", l)

	// Get len of range using "bytes=start-" syntax.
	l, err = httputil.LenOfRange(uri, 100000000, 0, true)
	if err != nil {
		log.Printf("httputil.LenOfRange() error: %v", err)
		return
	}

	fmt.Printf("size of range: 10000000-: %d\n", l)

	// Output:
	// size of range: 0 - 99999999: 100000000
	// size of range: 10000000-: 45565374
}

func ExampleLenOfRangeStart() {
	uri := "https://golang.google.cn/dl/go1.19.3.darwin-arm64.pkg"

	// Get len of range using "bytes=start-" syntax.
	l, err := httputil.LenOfRangeStart(uri, 100000000)
	if err != nil {
		log.Printf("httputil.LenOfRangeStart() error: %v", err)
		return
	}

	fmt.Printf("size of range: 10000000-: %d\n", l)

	// Output:
	// size of range: 10000000-: 45565374
}

func ExampleGetRespOfRange() {
	uri := "https://golang.google.cn/dl/go1.19.3.darwin-arm64.pkg"

	resp, l, err := httputil.GetRespOfRange(uri, 0, 99999999, false)
	if err != nil {
		log.Printf("httputil.GetRespOfRange() error: %v", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("size of range: 0 - 99999999: %d\n", l)

	// Get len of range using "bytes=start-" syntax.
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

	resp, l, err := httputil.GetRespOfRangeStart(uri, 100000000)
	if err != nil {
		log.Printf("httputil.LenOfRange() error: %v", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("size of range: 10000000-: %d\n", l)

	// Output:
	// size of range: 10000000-: 45565374
}
