package helper

import (
	"bytes"
	"io"
	"net/http"
)

func UploadImageHandler(req *http.Request) (res []byte, err error) {
	imageFile, _, err := req.FormFile("image")
	if err != nil {
		return
	}
	defer imageFile.Close()
	var imageBytes bytes.Buffer
	_, err = io.Copy(&imageBytes, imageFile)
	if err != nil {
		return
	}
	res = imageBytes.Bytes()
	return
}
