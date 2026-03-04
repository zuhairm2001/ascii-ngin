package video_test

import (
	"image"
	"image/color"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/zuhairm2001/ascii-ngin/pkg/video"
)

var _ = Describe("ReadPixelDump", func() {
	DescribeTable("reads pixel dump output",
		func(width int, height int, c color.RGBA, expectedRed int) {
			img := image.NewRGBA(image.Rect(0, 0, width, height))
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					img.SetRGBA(x, y, c)
				}
			}

			tmpDir := GinkgoT().TempDir()
			outputFile := filepath.Join(tmpDir, "pixel_dump.txt")
			err := video.WritePixelDump(img, outputFile)
			Expect(err).NotTo(HaveOccurred())

			pixels, err := video.ReadPixelDump(outputFile)
			Expect(err).NotTo(HaveOccurred())
			Expect(pixels).To(HaveLen(height))
			Expect(pixels[0]).To(HaveLen(width))
			Expect(pixels[0][0].Red).To(Equal(expectedRed))
		},
		Entry("small red image",
			2,
			2,
			color.RGBA{R: 100, G: 0, B: 0, A: 255},
			100,
		),
		Entry("small white image",
			1,
			1,
			color.RGBA{R: 255, G: 255, B: 255, A: 255},
			255,
		),
	)

	DescribeTable("surfaces parse errors",
		func(contents string) {
			tmpDir := GinkgoT().TempDir()
			outputFile := filepath.Join(tmpDir, "bad_dump.txt")
			err := os.WriteFile(outputFile, []byte(contents), 0644)
			Expect(err).NotTo(HaveOccurred())

			_, err = video.ReadPixelDump(outputFile)
			Expect(err).To(HaveOccurred())
		},
		Entry("invalid line format", "1,2: (not,a,pixel)\n"),
		Entry("missing coordinates", "(1,2,3)\n"),
	)
})
