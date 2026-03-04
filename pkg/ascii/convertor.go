package ascii

import (
	"log"
)

const (
	ASCII_CHARS         = " .'\\`^,:;Il!i><~+_-?][}{1)(|/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	LUMINANCE_THRESHOLD = 30
	FONT_RATIO          = 0.44
	LUMINANCE_RANGE     = 255.0
)

var logger = log.Default()

func CalculateLuminance(r, g, b int) float64 {
	return 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)
}

func MapLuminanceToASCII(luminance float64) rune {
	index := int((luminance / LUMINANCE_RANGE) * float64(len(ASCII_CHARS)-1))
	return rune(ASCII_CHARS[index])
}

func ScaledDimensions(imgWidth, imgHeight, termCols, termRows int) (cols, rows int) {
	logger.Println("scaling image with"+" imgWidth:", imgWidth, " imgHeight:", imgHeight, " termCols:", termCols, " termRows:", termRows)
	widthCols, widthRows := ScaleToWidth(imgWidth, imgHeight, termCols)

	if widthRows <= termRows {
		logger.Printf("scaling to width: %d x %d\n", widthCols, widthRows)
		return widthCols, widthRows
	}
	heightCols, heightRows := ScaleToHeight(imgWidth, imgHeight, termRows)
	logger.Printf("scaling to height: %d x %d\n", heightCols, heightRows)
	return heightCols, heightRows
	// return ScaleToHeight(imgWidth, imgHeight, termRows)
}

func ScaleToWidth(imgWidth, imgHeight, targetCols int) (cols, rows int) {
	cols = targetCols
	rows = int(float64(imgHeight) / float64(imgWidth) * float64(cols) * FONT_RATIO)
	return cols, rows
}

func ScaleToHeight(imgWidth, imgHeight, targetRows int) (cols, rows int) {
	rows = targetRows
	cols = int(float64(imgWidth) / float64(imgHeight) * float64(rows) / FONT_RATIO)
	return cols, rows
}
