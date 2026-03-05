package video_test

import (
	"image"
	"image/color"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/zuhairm2001/ascii-ngin/pkg/ascii"
	"github.com/zuhairm2001/ascii-ngin/pkg/video"
)

func TestImageToASCII(t *testing.T) {
	tests := []struct {
		name         string
		img          image.Image
		options      video.ASCIIOptions
		expectedRows int
		expectedCols int
		expectedRune rune
	}{
		{
			name:         "white pixel image with explicit size",
			img:          newSolidImage(2, 2, color.RGBA{R: 255, G: 255, B: 255, A: 255}),
			options:      video.ASCIIOptions{TermWidth: 2, TermHeight: 2},
			expectedRows: 1,
			expectedCols: 2,
			expectedRune: expectedRuneForColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}),
		},
		{
			name:         "black pixel image with explicit size",
			img:          newSolidImage(3, 3, color.RGBA{R: 0, G: 0, B: 0, A: 255}),
			options:      video.ASCIIOptions{TermWidth: 3, TermHeight: 3},
			expectedRows: 1,
			expectedCols: 3,
			expectedRune: expectedRuneForColor(color.RGBA{R: 0, G: 0, B: 0, A: 255}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			asciiArt, err := video.ImageToASCII(tt.img, "sample.png", tt.options)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(asciiArt).To(HaveLen(tt.expectedRows))
			g.Expect(asciiArt[0]).To(HaveLen(tt.expectedCols))
			for _, row := range asciiArt {
				for _, char := range row {
					g.Expect(char).To(Equal(tt.expectedRune))
				}
			}
		})
	}
}

func newSolidImage(width int, height int, c color.RGBA) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := range height {
		for x := range width {
			img.SetRGBA(x, y, c)
		}
	}
	return img
}

func expectedRuneForColor(c color.RGBA) rune {
	luminance := ascii.CalculateLuminance(int(c.R), int(c.G), int(c.B))
	if luminance < ascii.LUMINANCE_THRESHOLD {
		return ' '
	}
	return ascii.MapLuminanceToASCII(luminance)
}
