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
