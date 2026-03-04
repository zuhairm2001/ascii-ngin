package video

import (
	"image"

	"golang.org/x/image/draw"
)

func ResizeImage(src image.Image, newWidth int, newHeight int) image.Image {
	dstImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.BiLinear.Scale(dstImg, dstImg.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dstImg
}
