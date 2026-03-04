package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/zuhairm2001/ascii-ngin/pkg/video"
)

func main() {
	dumpDir := flag.String("dump-dir", "", "Write pixel dumps to directory")
	termSize := flag.String("term-size", "", "Override terminal size as WIDTHxHEIGHT or WIDTH,HEIGHT")
	flag.Parse()

	if err := video.CheckDependencies(); err != nil {
		return
	}

	termWidth, termHeight, err := parseTermSize(*termSize)
	if err != nil {
		fmt.Println("Invalid term-size:", err)
		return
	}

	imagePath := getImagePath("test.png")
	asciiArt, err := video.FileToASCII(imagePath, video.ASCIIOptions{
		TermWidth:    termWidth,
		TermHeight:   termHeight,
		PixelDumpDir: *dumpDir,
	})
	if err != nil {
		fmt.Println("Error generating ASCII art:", err)
		return
	}

	video.PrintASCIIArt(asciiArt)
	if err := video.WriteASCIIArtToFile(asciiArt, "output.txt"); err != nil {
		fmt.Println("Error writing ASCII art to file:", err)
		return
	}
	fmt.Println("ASCII Art generation completed.")
}

func getImagePath(filename string) string {
	baseDir := os.Getenv("ASCII_NGIN_IMAGES_DIR")
	if baseDir == "" {
		baseDir = "images"
	}
	return filepath.Join(baseDir, filename)
}

// we use this for testing purposes
// we can test various terminal sizes with this flag
func parseTermSize(value string) (int, int, error) {
	if strings.TrimSpace(value) == "" {
		return 0, 0, nil
	}

	separator := "x"
	if strings.Contains(value, ",") {
		separator = ","
	}

	parts := strings.Split(value, separator)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("expected WIDTHxHEIGHT or WIDTH,HEIGHT")
	}

	width, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || width <= 0 {
		return 0, 0, fmt.Errorf("invalid width")
	}

	height, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil || height <= 0 {
		return 0, 0, fmt.Errorf("invalid height")
	}

	return width, height, nil
}
