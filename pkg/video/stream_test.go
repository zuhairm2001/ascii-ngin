package video_test

import (
	"image"
	"image/color"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/zuhairm2001/ascii-ngin/pkg/video"
)

func TestASCIIFrameStream(t *testing.T) {
	tests := []struct {
		name           string
		images         []image.Image
		options        video.ASCIIOptions
		expectedCount  int
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "single image",
			images:         []image.Image{newSolidImage(2, 2, color.RGBA{R: 255, G: 255, B: 255, A: 255})},
			options:        video.ASCIIOptions{TermWidth: 2, TermHeight: 2},
			expectedCount:  1,
			expectedWidth:  2,
			expectedHeight: 1,
		},
		{
			name: "multiple images",
			images: []image.Image{
				newSolidImage(3, 3, color.RGBA{R: 0, G: 0, B: 0, A: 255}),
				newSolidImage(3, 3, color.RGBA{R: 0, G: 0, B: 0, A: 255}),
			},
			options:        video.ASCIIOptions{TermWidth: 3, TermHeight: 3},
			expectedCount:  2,
			expectedWidth:  3,
			expectedHeight: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			imageChan := make(chan image.Image, len(tt.images))
			for _, img := range tt.images {
				imageChan <- img
			}
			close(imageChan)

			frames, errs := video.ASCIIFrameStream(imageChan, nil, tt.options)

			var collected []video.FrameRecord
			for frame := range frames {
				collected = append(collected, frame)
			}

			for err := range errs {
				g.Expect(err).NotTo(HaveOccurred())
			}

			g.Expect(collected).To(HaveLen(tt.expectedCount))
			for index, frame := range collected {
				g.Expect(frame.I).To(Equal(index))
				g.Expect(frame.W).To(Equal(tt.expectedWidth))
				g.Expect(frame.H).To(Equal(tt.expectedHeight))
				g.Expect(frame.F).NotTo(BeEmpty())
			}
		})
	}
}
