package video_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestVideo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Video Suite")
}
