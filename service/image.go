package service

import (
	"image"
	"image/jpeg"
	"io"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"github.com/workdestiny/watgok_web/config"
)

//ResizeDisplay resize image
func ResizeDisplay(m image.Image) image.Image {
	return resize.Resize(config.WidthDisplay, config.HeightDisplay, m, resize.Lanczos3)
}

//GenerateDisplayName new url path
func GenerateDisplayName(id string) string {
	return "profile/display/" + id + "-" + uuid.New().String()
}

//EncodeJPEG encode
func EncodeJPEG(w io.Writer, m image.Image, quality int) error {
	return jpeg.Encode(w, m, &jpeg.Options{Quality: quality})
}
