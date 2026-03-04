// CLI usage:
//
//	go run ./cmd/ascii-ngin [flags] <video>
//
// Flags:
//
//	--fps=<int>           Frames per second for playback (default: 24)
//	--term-size=WxH       Override terminal size (example: 120x40 or 120,40)
//	--dump-dir=<path>     Write pixel dumps to directory
//	--ndjson-out=<path>   Write NDJSON frames to file
//	--ndjson-in=<path>    Read NDJSON frames from file
//	--debug               Enable debug logging
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/zuhairm2001/ascii-ngin/pkg/ascii"
	"github.com/zuhairm2001/ascii-ngin/pkg/video"
)

func main() {
	dumpDir := flag.String("dump-dir", "", "Write pixel dumps to directory")
	ndjsonOut := flag.String("ndjson-out", "", "Write NDJSON frames to file")
	ndjsonIn := flag.String("ndjson-in", "", "Read NDJSON frames from file")
	fps := flag.Int("fps", 24, "Frames per second for playback")
	debug := flag.Bool("debug", false, "Enable debug logging")
	termSize := flag.String("term-size", "", "Override terminal size as WIDTHxHEIGHT or WIDTH,HEIGHT")
	flag.Parse()

	if err := video.CheckDependencies(); err != nil {
		return
	}

	ascii.Debug = *debug

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
	if *fps <= 0 {
		fmt.Println("Invalid fps:", *fps)
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage: ascii-ngin [flags] <video>")
		return
	}

	videoPath := args[0]
	images, imageErrs := video.ExtractFramesPipe(videoPath, *fps)
	frames, frameErrs := video.ASCIIFrameStream(images, nil, video.ASCIIOptions{
		TermWidth:    termWidth,
		TermHeight:   termHeight,
		PixelDumpDir: *dumpDir,
	})

	if err := renderStream(videoPath, frames, frameErrs, imageErrs, *fps, *ndjsonOut); err != nil {
		fmt.Println("Error rendering stream:", err)
		return
	}
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
		clearScreen()
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

func renderStream(source string, frames <-chan video.FrameRecord, frameErrs <-chan error, imageErrs <-chan error, fps int, ndjsonOut string) error {
	interval := time.Second / time.Duration(fps)

	firstFrame, ok := <-frames
	if !ok {
		return drainErrors(frameErrs, imageErrs)
	}

	writer, encoder, closeWriter, err := openNDJSONWriter(ndjsonOut, video.MetaRecord{
		Type:   "meta",
		FPS:    fps,
		W:      firstFrame.W,
		H:      firstFrame.H,
		Source: source,
	})
	if err != nil {
		return err
	}
	if closeWriter != nil {
		defer closeWriter()
	}

	if encoder != nil {
		if err := encoder.Encode(firstFrame); err != nil {
			return err
		}
	}
	clearScreen()
	if _, err := io.WriteString(os.Stdout, firstFrame.F+"\n"); err != nil {
		return err
	}
	time.Sleep(interval)

	for frame := range frames {
		if encoder != nil {
			if err := encoder.Encode(frame); err != nil {
				return err
			}
		}
		clearScreen()
		if _, err := io.WriteString(os.Stdout, frame.F+"\n"); err != nil {
			return err
		}
		time.Sleep(interval)
	}

	if err := drainErrors(frameErrs, imageErrs); err != nil {
		return err
	}

	if writer != nil {
		if err := writer.Sync(); err != nil {
			return err
		}
	}

	return nil
}

func openNDJSONWriter(path string, meta video.MetaRecord) (*os.File, *json.Encoder, func(), error) {
	if path == "" {
		return nil, nil, nil, nil
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, nil, nil, err
	}

	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(meta); err != nil {
		_ = file.Close()
		return nil, nil, nil, err
	}

	return file, encoder, func() { _ = file.Close() }, nil
}

func drainErrors(frameErrs <-chan error, imageErrs <-chan error) error {
	for err := range frameErrs {
		if err != nil {
			return err
		}
	}
	for err := range imageErrs {
		if err != nil {
			return err
		}
	}
	return nil
}

func clearScreen() {
	_, _ = io.WriteString(os.Stdout, "\x1b[2J\x1b[H")
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
