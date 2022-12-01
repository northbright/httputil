package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/northbright/pathelper"
)

var (
	serverRoot = ""
)

func init() {
	serverRoot, _ = pathelper.ExecDir("")
}

func uploadFile(c *gin.Context) {
	var err error

	defer func() {
		if err != nil {
			log.Printf("%v", err)
		}
	}()

	file, header, err := c.Request.FormFile("upload")
	if err != nil {
		err = fmt.Errorf("FormFile() error: %v", err)
		return
	}

	uploadedDir := path.Join(serverRoot, "uploaded")
	fileName := path.Join(uploadedDir, header.Filename)

	out, err := os.Create(fileName)
	if err != nil {
		err = fmt.Errorf("os.Create() error: %v", err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		err = fmt.Errorf("io.Copy() error: %v", err)
		return
	}
}

func main() {
	var err error

	defer func() {
		if err != nil {
			log.Printf("main() err: %v\n", err)
		}
	}()

	r := gin.Default()

	r.POST("/", uploadFile)

	r.Run(":80")
}
