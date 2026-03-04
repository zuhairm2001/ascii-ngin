package video

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"os/exec"
)

func ExtractFramesPipe(videoPath string, fps int) (<-chan image.Image, <-chan error) {
	images := make(chan image.Image)
	errs := make(chan error, 1)

	go func() {
		defer close(images)
		defer close(errs)

		if fps <= 0 {
			errs <- fmt.Errorf("fps must be positive")
			return
		}

		cmd := exec.Command("ffmpeg", "-i", videoPath, "-vf", fmt.Sprintf("fps=%d", fps), "-f", "image2pipe", "-vcodec", "png", "-")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			errs <- err
			return
		}

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Start(); err != nil {
			errs <- err
			return
		}

		reader := bufio.NewReader(stdout)
		for {
			img, _, err := image.Decode(reader)
			if err != nil {
				if err == io.EOF {
					break
				}
				errs <- fmt.Errorf("failed to decode frame: %w", err)
				_ = cmd.Wait()
				return
			}
			images <- img
		}

		if err := cmd.Wait(); err != nil {
			errs <- fmt.Errorf("ffmpeg error: %w: %s", err, stderr.String())
			return
		}
	}()

	return images, errs
}
