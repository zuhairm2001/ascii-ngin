package video

import (
	"bufio"
	"fmt"
	"image"
	"os"
	"path/filepath"
)

func WritePixelDump(img image.Image, outputFile string) error {
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	bounds := img.Bounds()
	fmt.Fprintf(outFile, "# ImageMagick pixel enumeration: %d,%d,255,srgb\n", bounds.Dx(), bounds.Dy())

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			fmt.Fprintf(outFile, "%d,%d: (%d,%d,%d)\n", x, y, r>>8, g>>8, b>>8)
		}
	}

	return nil
}

func ReadPixelDump(filePath string) ([][]Pixel, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var pixelData [][]Pixel
	var currentRowPixelArr []Pixel
	currentRow := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] == '#' {
			continue
		}

		pixel, err := parsePixelLine(line)
		if err != nil {
			return nil, err
		}

		if pixel.Y != currentRow {
			if len(currentRowPixelArr) > 0 {
				pixelData = append(pixelData, currentRowPixelArr)
			}
			currentRowPixelArr = []Pixel{}
			currentRow = pixel.Y
		}

		currentRowPixelArr = append(currentRowPixelArr, pixel)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(currentRowPixelArr) > 0 {
		pixelData = append(pixelData, currentRowPixelArr)
	}

	return pixelData, nil
}

func parsePixelLine(line string) (Pixel, error) {
	var x, y, r, g, b int
	_, err := fmt.Sscanf(line, "%d,%d: (%d,%d,%d)", &x, &y, &r, &g, &b)
	if err != nil {
		return Pixel{}, err
	}
	return Pixel{X: x, Y: y, Red: r, Green: g, Blue: b}, nil
}
