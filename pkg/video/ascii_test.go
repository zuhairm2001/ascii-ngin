package video_test

import (
	"image"
	"image/color"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/zuhairm2001/ascii-ngin/pkg/ascii"
	"github.com/zuhairm2001/ascii-ngin/pkg/video"
)

var _ = Describe("ImageToASCII", func() {
	DescribeTable("creates ASCII art from images",
		func(img image.Image, options video.ASCIIOptions, expectedRows int, expectedCols int, expectedRune rune) {
			asciiArt, err := video.ImageToASCII(img, "sample.png", options)
			Expect(err).NotTo(HaveOccurred())
			Expect(asciiArt).To(HaveLen(expectedRows))
			Expect(asciiArt[0]).To(HaveLen(expectedCols))
			for _, row := range asciiArt {
				for _, char := range row {
					Expect(char).To(Equal(expectedRune))
				}
			}
		},
		Entry("white pixel image with explicit size",
			newSolidImage(2, 2, color.RGBA{R: 255, G: 255, B: 255, A: 255}),
			video.ASCIIOptions{TermWidth: 2, TermHeight: 2},
			1,
			2,
			expectedRuneForColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}),
		),
		Entry("black pixel image with explicit size",
			newSolidImage(3, 3, color.RGBA{R: 0, G: 0, B: 0, A: 255}),
			video.ASCIIOptions{TermWidth: 3, TermHeight: 3},
			1,
			3,
			expectedRuneForColor(color.RGBA{R: 0, G: 0, B: 0, A: 255}),
		),
	)
})

func newSolidImage(width int, height int, c color.RGBA) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
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
