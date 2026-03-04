package video

import (
	"fmt"
	"image"
)

const defaultFrameNameFormat = "frame_%06d.png"

func ASCIIFrameStream(images <-chan image.Image, nameForIndex func(int) string, options ASCIIOptions) (<-chan FrameRecord, <-chan error) {
	frames := make(chan FrameRecord, 1)
	errs := make(chan error, 1)

	go func() {
		defer close(frames)
		defer close(errs)

		frameIndex := 0
		for img := range images {
			name := frameName(nameForIndex, frameIndex)
			asciiArt, err := ImageToASCII(img, name, options)
			if err != nil {
				errs <- err
				return
			}

			frame := FrameRecord{
				I: frameIndex,
				W: frameWidth(asciiArt),
				H: len(asciiArt),
				F: ASCIIArtToString(asciiArt),
			}

			frames <- frame
			frameIndex++
		}
	}()

	return frames, errs
}

func frameName(nameForIndex func(int) string, index int) string {
	if nameForIndex != nil {
		return nameForIndex(index)
	}
	return fmt.Sprintf(defaultFrameNameFormat, index)
}

func frameWidth(asciiArt [][]rune) int {
	if len(asciiArt) == 0 {
		return 0
	}
	return len(asciiArt[0])
}
