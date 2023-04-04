package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/northbright/pathelper"
)

var (
	serverRoot = ""
)

func init() {
	serverRoot, _ = pathelper.ExecDir("")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file_to_upload")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	uploadedDir := path.Join(serverRoot, "uploaded")
	fileName := path.Join(uploadedDir, header.Filename)

	out, err := os.Create(fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("File is uploaded successfully."))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", uploadHandler)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
