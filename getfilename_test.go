package httputil_test

import (
	"fmt"
	"os"

	"github.com/northbright/httputil"
)

func ExampleGetFileName() {
	url := "https://github.com/northbright/plants/archive/master.zip"

	if f, err := httputil.GetFileName(url); err != nil {
		fmt.Fprintf(os.Stderr, "GetFileName(%v, resp) err: %v\n", url, err)
	} else {
		fmt.Fprintf(os.Stderr, "GetFileName() succeeded: file name = %v\n", f)
	}
	// Output:
}
