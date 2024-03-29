package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/northbright/httputil"
)

func main() {
	var (
		err        error
		uri        = "http://localhost:8080"
		uploadFile = "files/1.txt"
		client     = http.Client{}
		req        *http.Request
		resp       *http.Response
		data       []byte
	)

	defer func() {
		if err != nil {
			log.Printf("%v", err)
		}
	}()

	execDir, _ := os.Executable()
	uploadFile = filepath.Join(execDir, uploadFile)
	if req, err = httputil.NewUploadFileRequest(
		"POST",
		uri,
		uploadFile,
		"file_to_upload",
		nil); err != nil {
		err = fmt.Errorf("NewUploadFileRequest() error: %v", err)
		return
	}

	if resp, err = client.Do(req); err != nil {
		err = fmt.Errorf("client.Do() error: %v", err)
		return
	}
	defer resp.Body.Close()

	if data, err = ioutil.ReadAll(resp.Body); err != nil {
		err = fmt.Errorf("ioutil.ReadAll() error: %v", err)
		return
	}

	log.Printf("Response StatusCode: %v\nBody: %v\n", resp.StatusCode, string(data))
}
