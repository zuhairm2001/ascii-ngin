package video

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

func GetImageDimensions(filename string) (int, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to decode image config: %w", err)
	}

	return config.Height, config.Width, nil
}

func ResizeImage(inputFile string, outputFile string, newWidth int, newHeight int) error {
	srcFile, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer srcFile.Close()

	srcImg, _, err := image.Decode(srcFile)
	if err != nil {
		return fmt.Errorf("failed to decode input image: %w", err)
	}

	dstImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.BiLinear.Scale(dstImg, dstImg.Bounds(), srcImg, srcImg.Bounds(), draw.Over, nil)

	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	ext := strings.ToLower(filepath.Ext(outputFile))
	switch ext {
	case ".png":
		if err := png.Encode(outFile, dstImg); err != nil {
			return fmt.Errorf("failed to encode PNG: %w", err)
		}
	case ".jpg", ".jpeg":
		if err := jpeg.Encode(outFile, dstImg, &jpeg.Options{Quality: 90}); err != nil {
			return fmt.Errorf("failed to encode JPEG: %w", err)
		}
	default:
		if err := jpeg.Encode(outFile, dstImg, &jpeg.Options{Quality: 90}); err != nil {
			return fmt.Errorf("failed to encode image: %w", err)
		}
	}

	return nil
}

func ConvertToTextPixelData(imageFile string, outputFile string) error {
	srcFile, err := os.Open(imageFile)
	if err != nil {
		return fmt.Errorf("failed to open image file: %w", err)
	}
	defer srcFile.Close()

	img, _, err := image.Decode(srcFile)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	bounds := img.Bounds()
	fmt.Fprintf(outFile, "# ImageMagick pixel enumeration: %d,%d,255,srgb\n",
		bounds.Dx(), bounds.Dy())

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// RGBA() returns 16-bit values, convert to 8-bit (0-255)
			fmt.Fprintf(outFile, "%d,%d: (%d,%d,%d)\n", x, y, r>>8, g>>8, b>>8)
		}
	}

	return nil
}
