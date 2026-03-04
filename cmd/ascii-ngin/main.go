package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/zuhairm2001/ascii-ngin/pkg/video"
)

func main() {
	dumpDir := flag.String("dump-dir", "", "Write pixel dumps to directory")
	ndjsonOut := flag.String("ndjson-out", "", "Write NDJSON frames to file")
	ndjsonIn := flag.String("ndjson-in", "", "Read NDJSON frames from file")
	fps := flag.Int("fps", 24, "Frames per second for playback")
	termSize := flag.String("term-size", "", "Override terminal size as WIDTHxHEIGHT or WIDTH,HEIGHT")
	flag.Parse()

	if err := video.CheckDependencies(); err != nil {
		return
	}

	if *ndjsonIn != "" {
		if err := playNDJSON(*ndjsonIn, *fps); err != nil {
			fmt.Println("Error playing NDJSON:", err)
		}
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

	if *ndjsonOut != "" {
		if err := writeNDJSONOutput(*ndjsonOut, imagePath, *fps, asciiArt); err != nil {
			fmt.Println("Error writing NDJSON:", err)
			return
		}
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

func writeNDJSONOutput(path string, source string, fps int, asciiArt [][]rune) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	frames := make(chan video.FrameRecord, 1)
	frames <- video.FrameRecord{
		I: 0,
		W: frameWidth(asciiArt),
		H: len(asciiArt),
		F: video.ASCIIArtToString(asciiArt),
	}
	close(frames)

	meta := video.MetaRecord{
		Type:   "meta",
		FPS:    fps,
		W:      frameWidth(asciiArt),
		H:      len(asciiArt),
		Source: source,
	}

	return video.WriteNDJSON(frames, file, meta)
}

func playNDJSON(path string, fps int) error {
	if fps <= 0 {
		return fmt.Errorf("fps must be positive")
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	frames, errs := video.ReadNDJSON(file)
	interval := time.Second / time.Duration(fps)
	for frame := range frames {
		if _, err := io.WriteString(os.Stdout, frame.F+"\n"); err != nil {
			return err
		}
		time.Sleep(interval)
	}

	if err, ok := <-errs; ok && err != nil {
		return err
	}

	return nil
}

func frameWidth(asciiArt [][]rune) int {
	if len(asciiArt) == 0 {
		return 0
	}
	return len(asciiArt[0])
}
