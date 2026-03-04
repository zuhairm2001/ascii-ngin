package video_test

import (
	"bytes"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/zuhairm2001/ascii-ngin/pkg/video"
)

var _ = Describe("NDJSON", func() {
	DescribeTable("writes and reads frames",
		func(frames []video.FrameRecord, meta video.MetaRecord) {
			var buffer bytes.Buffer
			frameChan := make(chan video.FrameRecord, len(frames))
			for _, frame := range frames {
				frameChan <- frame
			}
			close(frameChan)

			err := video.WriteNDJSON(frameChan, &buffer, meta)
			Expect(err).NotTo(HaveOccurred())

			readFrames, errs := video.ReadNDJSON(strings.NewReader(buffer.String()))
			var collected []video.FrameRecord
			for frame := range readFrames {
				collected = append(collected, frame)
			}
			Expect(<-errs).NotTo(HaveOccurred())
			Expect(collected).To(Equal(frames))
		},
		Entry("single frame",
			[]video.FrameRecord{{I: 0, W: 2, H: 1, F: "@@"}},
			video.MetaRecord{Type: "meta", FPS: 24, W: 2, H: 1, Source: "test.mp4"},
		),
		Entry("multiple frames",
			[]video.FrameRecord{
				{I: 0, W: 2, H: 1, F: "@@"},
				{I: 1, W: 2, H: 1, F: "##"},
			},
			video.MetaRecord{FPS: 30, W: 2, H: 1},
		),
	)

	DescribeTable("returns errors for invalid ndjson",
		func(contents string) {
			readFrames, errs := video.ReadNDJSON(strings.NewReader(contents))
			for range readFrames {
			}
			Expect(<-errs).To(HaveOccurred())
		},
		Entry("invalid json", "{not-json}\n"),
		Entry("partial json", "{\"i\":0\n"),
	)
})
