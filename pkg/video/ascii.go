package video

import (
	"bufio"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/zuhairm2001/ascii-ngin/pkg/ascii"
)

type ASCIIOptions struct {
	TermWidth    int
	TermHeight   int
	PixelDumpDir string
}

func FileToASCII(filename string, options ASCIIOptions) ([][]rune, error) {
	img, err := DecodeImage(filename)
	if err != nil {
		return nil, err
	}

	if options.TermWidth == 0 || options.TermHeight == 0 {
		termData, err := GetTerminalDimensions()
		if err != nil {
			return nil, err
		}
		if options.TermWidth == 0 {
			options.TermWidth = termData.Width
		}
		if options.TermHeight == 0 {
			options.TermHeight = termData.Height
		}
	}

	return ImageToASCII(img, filepath.Base(filename), options)
}

func ImageToASCII(img image.Image, name string, options ASCIIOptions) ([][]rune, error) {
	if options.TermWidth == 0 || options.TermHeight == 0 {
		termData, err := GetTerminalDimensions()
		if err == nil {
			if options.TermWidth == 0 {
				options.TermWidth = termData.Width
			}
			if options.TermHeight == 0 {
				options.TermHeight = termData.Height
			}
		}
	}
	if options.TermWidth == 0 {
		options.TermWidth = 80
	}
	if options.TermHeight == 0 {
		options.TermHeight = 24
	}

	bounds := img.Bounds()
	newWidth, newHeight := ascii.ScaledDimensions(bounds.Dx(), bounds.Dy(), options.TermWidth, options.TermHeight)
	if newWidth < 1 {
		newWidth = 1
	}
	if newHeight < 1 {
		newHeight = 1
	}
	resized := ResizeImage(img, newWidth, newHeight)

	if options.PixelDumpDir != "" {
		dumpPath := filepath.Join(options.PixelDumpDir, "text_"+name+".txt")
		if err := WritePixelDump(resized, dumpPath); err != nil {
			return nil, err
		}
	}

	var asciiArt [][]rune
	for y := 0; y < resized.Bounds().Dy(); y++ {
		var asciiRow []rune
		for x := 0; x < resized.Bounds().Dx(); x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			asciiChar := PixelToASCII(Pixel{Red: int(r >> 8), Green: int(g >> 8), Blue: int(b >> 8)})
			asciiRow = append(asciiRow, asciiChar)
		}
		asciiArt = append(asciiArt, asciiRow)
	}

	return asciiArt, nil
}

// given rbg values of a pixel return the corresponding ascii character
func PixelToASCII(pixel Pixel) rune {
	luminance := ascii.CalculateLuminance(pixel.Red, pixel.Green, pixel.Blue)

	if luminance < ascii.LUMINANCE_THRESHOLD {
		return ' '
	}

	return ascii.MapLuminanceToASCII(luminance)
}

func PrintASCIIArt(asciiArt [][]rune) {
	for _, row := range asciiArt {
		for _, char := range row {
			fmt.Print(string(char))
		}
		fmt.Println()
	}
}

func WriteASCIIArtToFile(asciiArt [][]rune, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, row := range asciiArt {
		for _, char := range row {
			_, err := writer.WriteString(string(char))
			if err != nil {
				return err
			}
		}
		_, err := writer.WriteString("\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func ASCIIArtToString(asciiArt [][]rune) string {
	if len(asciiArt) == 0 {
		return ""
	}

	var builder strings.Builder
	for rowIndex, row := range asciiArt {
		for _, char := range row {
			builder.WriteRune(char)
		}
		if rowIndex < len(asciiArt)-1 {
			builder.WriteByte('\n')
		}
	}

	return builder.String()
}
