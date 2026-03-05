package video_test

import (
	"bytes"
	"strings"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/zuhairm2001/ascii-ngin/pkg/video"
)

func TestNDJSON_WriteAndReadFrames(t *testing.T) {
	tests := []struct {
		name   string
		frames []video.FrameRecord
		meta   video.MetaRecord
	}{
		{
			name:   "single frame",
			frames: []video.FrameRecord{{I: 0, W: 2, H: 1, F: "@@"}},
			meta:   video.MetaRecord{Type: "meta", FPS: 24, W: 2, H: 1, Source: "test.mp4"},
		},
		{
			name: "multiple frames",
			frames: []video.FrameRecord{
				{I: 0, W: 2, H: 1, F: "@@"},
				{I: 1, W: 2, H: 1, F: "##"},
			},
			meta: video.MetaRecord{FPS: 30, W: 2, H: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			var buffer bytes.Buffer
			frameChan := make(chan video.FrameRecord, len(tt.frames))
			for _, frame := range tt.frames {
				frameChan <- frame
			}
			close(frameChan)

			err := video.WriteNDJSON(frameChan, &buffer, tt.meta)
			g.Expect(err).NotTo(HaveOccurred())

			readFrames, errs := video.ReadNDJSON(strings.NewReader(buffer.String()))
			var collected []video.FrameRecord
			for frame := range readFrames {
				collected = append(collected, frame)
			}
			g.Expect(<-errs).NotTo(HaveOccurred())
			g.Expect(collected).To(Equal(tt.frames))
		})
	}
}

func TestNDJSON_InvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		contents string
	}{
		{name: "invalid json", contents: "{not-json}\n"},
		{name: "partial json", contents: "{\"i\":0\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			readFrames, errs := video.ReadNDJSON(strings.NewReader(tt.contents))
			for range readFrames {
			}
			g.Expect(<-errs).To(HaveOccurred())
		})
	}
}
