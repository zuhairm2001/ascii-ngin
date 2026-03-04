package video_test

import (
	"image"
	"image/color"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/zuhairm2001/ascii-ngin/pkg/video"
)

var _ = Describe("ASCIIFrameStream", func() {
	DescribeTable("converts images into frame records",
		func(images []image.Image, options video.ASCIIOptions, expectedCount int, expectedWidth int, expectedHeight int) {
			imageChan := make(chan image.Image, len(images))
			for _, img := range images {
				imageChan <- img
			}
			close(imageChan)

			frames, errs := video.ASCIIFrameStream(imageChan, nil, options)

			var collected []video.FrameRecord
			for frame := range frames {
				collected = append(collected, frame)
			}

			for err := range errs {
				Expect(err).NotTo(HaveOccurred())
			}

			Expect(collected).To(HaveLen(expectedCount))
			for index, frame := range collected {
				Expect(frame.I).To(Equal(index))
				Expect(frame.W).To(Equal(expectedWidth))
				Expect(frame.H).To(Equal(expectedHeight))
				Expect(frame.F).NotTo(BeEmpty())
			}
		},
		Entry("single image",
			[]image.Image{newSolidImage(2, 2, color.RGBA{R: 255, G: 255, B: 255, A: 255})},
			video.ASCIIOptions{TermWidth: 2, TermHeight: 2},
			1,
			2,
			1,
		),
		Entry("multiple images",
			[]image.Image{
				newSolidImage(3, 3, color.RGBA{R: 0, G: 0, B: 0, A: 255}),
				newSolidImage(3, 3, color.RGBA{R: 0, G: 0, B: 0, A: 255}),
			},
			video.ASCIIOptions{TermWidth: 3, TermHeight: 3},
			2,
			3,
			1,
		),
	)
})
