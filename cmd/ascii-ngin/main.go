package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/zuhairm2001/ascii-ngin/pkg/video"
)

func main() {
	ndjsonOut := flag.String("ndjson-out", "", "Write NDJSON frames to file")
	ndjsonIn := flag.String("ndjson-in", "", "Read NDJSON frames from file")
	fps := flag.Int("fps", 24, "Frames per second for playback")
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

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage: ascii-ngin [flags] <video>")
		return
	}

	videoPath := args[0]
	images, errs := video.ExtractFramesPipe(videoPath, *fps)
	for range images {
	}

	for err := range errs {
		if err != nil {
			fmt.Println("Error extracting frames:", err)
			return
		}
	}

	if *ndjsonOut != "" {
		fmt.Println("NDJSON output is not wired yet:", *ndjsonOut)
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
