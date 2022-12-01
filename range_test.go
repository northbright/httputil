package httputil_test

import (
	"fmt"
	"log"

	"github.com/northbright/httputil"
)

func ExampleIsRangeSupported() {
	url := "https://golang.google.cn/dl/go1.19.3.darwin-arm64.pkg"

	supported, err := httputil.IsRangeSupported(url)
	if err != nil {
		log.Printf("IsRangeSupported() error: %v", err)
		return
	}

	fmt.Print(supported)

	// Output:
	// true
}
