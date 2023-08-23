package helper

import (
	"bytes"
	"errors"
	"image"
	"io"
	"net/http"

	_ "image/jpeg"
	_ "image/png"

	"github.com/rs/zerolog/log"
)

var (
	imgHeigh = 400
	imgWidth = 600
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

	m, _, err := image.Decode(bytes.NewReader(imageBytes.Bytes()))
	if err != nil {
		log.Error().Err(err)
	}

	bounds := m.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	if w != imgWidth || h != imgHeigh {
		err = errors.New("image dimension must 600x400")
		return
	}

	res = imageBytes.Bytes()
	return
}
