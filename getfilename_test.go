package httputil_test

import (
	"fmt"
	"net/http"
	"os"

	"github.com/northbright/httputil"
)

func ExampleGetFileName() {
	url := "https://github.com/northbright/plants/archive/master.zip"
	if resp, err := http.Head(url); err != nil {
		fmt.Fprintf(os.Stderr, "http.Head(%v) err: %v\n", url, err)
	} else {
		defer resp.Body.Close()

		if f, err := httputil.GetFileName(url, resp); err != nil {
			fmt.Fprintf(os.Stderr, "GetFileName(%v, resp) err: %v\n", url, err)
		} else {
			fmt.Fprintf(os.Stderr, "GetFileName() succeeded: file name = %v\n", f)
		}
	}
	// Output:
}
