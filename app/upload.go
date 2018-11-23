package app

import (
	"context"
	"image"

	"github.com/workdestiny/watgok_web/service"
)

func upload(ctx context.Context, m image.Image, filename string) error {

	writer := bucket.Storage.Object(filename).NewWriter(ctx)
	writer.CacheControl = "public, max-age=31536000"
	defer writer.Close()
	err := service.EncodeJPEG(writer, m, 100)
	if err != nil {
		return err
	}
	return nil
}
