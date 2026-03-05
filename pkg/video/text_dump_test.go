package video_test

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/zuhairm2001/ascii-ngin/pkg/video"
)

func TestReadPixelDump(t *testing.T) {
	tests := []struct {
		name        string
		width       int
		height      int
		c           color.RGBA
		expectedRed int
	}{
		{
			name:        "small red image",
			width:       2,
			height:      2,
			c:           color.RGBA{R: 100, G: 0, B: 0, A: 255},
			expectedRed: 100,
		},
		{
			name:        "small white image",
			width:       1,
			height:      1,
			c:           color.RGBA{R: 255, G: 255, B: 255, A: 255},
			expectedRed: 255,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			img := image.NewRGBA(image.Rect(0, 0, tt.width, tt.height))
			for y := 0; y < tt.height; y++ {
				for x := 0; x < tt.width; x++ {
					img.SetRGBA(x, y, tt.c)
				}
			}

			tmpDir := t.TempDir()
			outputFile := filepath.Join(tmpDir, "pixel_dump.txt")
			err := video.WritePixelDump(img, outputFile)
			g.Expect(err).NotTo(HaveOccurred())

			pixels, err := video.ReadPixelDump(outputFile)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(pixels).To(HaveLen(tt.height))
			g.Expect(pixels[0]).To(HaveLen(tt.width))
			g.Expect(pixels[0][0].Red).To(Equal(tt.expectedRed))
		})
	}
}

func TestReadPixelDump_ParseErrors(t *testing.T) {
	tests := []struct {
		name     string
		contents string
	}{
		{name: "invalid line format", contents: "1,2: (not,a,pixel)\n"},
		{name: "missing coordinates", contents: "(1,2,3)\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			tmpDir := t.TempDir()
			outputFile := filepath.Join(tmpDir, "bad_dump.txt")
			err := os.WriteFile(outputFile, []byte(tt.contents), 0644)
			g.Expect(err).NotTo(HaveOccurred())

			_, err = video.ReadPixelDump(outputFile)
			g.Expect(err).To(HaveOccurred())
		})
	}
}
