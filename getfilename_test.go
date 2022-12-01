package httputil_test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/northbright/httputil"
)

func ExampleGetFileNameFromURL() {
	url := "https://github.com/northbright/plants/archive/master.zip"

	fileName, _ := httputil.GetFileNameFromURL(url)
	fmt.Print(fileName)

	// Output:
	// master.zip
}

func ExampleGetFileNameFromResponse() {
	url := "https://github.com/northbright/plants/archive/master.zip"

	resp, err := http.Head(url)
	if err != nil {
		log.Printf("head error: %v", err)
		return
	}
	defer resp.Body.Close()

	fileName, _ := httputil.GetFileNameFromResponse(resp)
	fmt.Print(fileName)

	// Output:
	// plants-master.zip
}
