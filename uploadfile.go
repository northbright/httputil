package httputil

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// NewUploadFileRequest() creates HTTP request to upload file.
//
// Params:
//     method: HTTP request method. Available values: "POST", "PUT".
//     uri: Request URL.
//     filePath: Absolute path of file to be uploaded.
//     fieldName: Field name(param name) of the uploaded file. Server side need this name to get uploaded file.
//     params: Extra params to be set in the form.
func NewUploadFileRequest(method, uri, filePath, fieldName string, params map[string]string) (*http.Request, error) {
	var (
		err  error
		f    *os.File
		body = &bytes.Buffer{}
		w    *multipart.Writer
		part io.Writer
		req  *http.Request
	)

	if method != "POST" && method != "PUT" {
		method = "POST"
	}

	if f, err = os.Open(filePath); err != nil {
		return nil, err
	}
	defer f.Close()

	w = multipart.NewWriter(body)

	if part, err = w.CreateFormFile(fieldName, filepath.Base(filePath)); err != nil {
		return nil, err
	}

	if _, err = io.Copy(part, f); err != nil {
		return nil, err
	}

	for k, v := range params {
		if err = w.WriteField(k, v); err != nil {
			return nil, err
		}
	}

	if err = w.Close(); err != nil {
		return nil, err
	}

	if req, err = http.NewRequest(method, uri, body); err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	return req, nil
}
