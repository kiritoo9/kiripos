package helpers

import (
	"encoding/base64"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func GenerateImage(dir string, base64Image string, unique string) string {
	var filename string
	if base64Image != "" {
		var path string = "./cdn/"

		// CREATE "CDN" FOLDER
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.Mkdir(path, 0700); err != nil {
				panic(err.Error())
			}
		}

		// CREATE SPECIFIC FOLDER
		path += dir + "/"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.Mkdir(path, 0700); err != nil {
				panic(err.Error())
			}
		}

		filename = unique + ".png"
		b64data := base64Image[strings.IndexByte(base64Image, ',')+1:]
		dec, err := base64.StdEncoding.DecodeString(b64data)
		if err != nil {
			panic(err.Error())
		}
		f, err := os.Create(path + filename)
		if err != nil {
			panic(err.Error())
		}
		defer f.Close()
		if _, err := f.Write(dec); err != nil {
			panic(err.Error())
		}
		if err := f.Sync(); err != nil {
			panic(err.Error())
		}
	}
	return filename
}

func GettRealPath(c *gin.Context, path string) string {
	var protocol string = "http"
	if c.Request.TLS != nil {
		protocol = "https"
	}
	return protocol + "://" + c.Request.Host + "/cdn/" + path
}
