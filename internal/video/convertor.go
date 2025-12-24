package video

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zuhairm2001/ascii-ngin/pkg/ascii"
)

type PixelData struct {
	Red   int
	Green int
	Blue  int
}

type FrameData struct {
	Width    int
	Height   int
	FrameNum int
}

type VideoData struct {
	FrameCount int
	FrameRate  float64
	Frames     []FrameData
	Metadata   VideoMetadata
}

type VideoMetadata struct {
	Duration   float64
	Resolution string
	FPS        float64
	Filename   string
}

type Video interface {
	GetFrameCount() int
	GetFrameRate() float64
	GetFrames() []FrameData
}

func Run() {
	err := CheckDependencies()
	if err != nil {
		return
	}
	fmt.Println("All dependencies are installed. Proceeding with video processing...")
	imagePath := getImagePath("test.png")
	fmt.Print(imagePath)
	asciiArt := FrameToASCII(FrameData{}, imagePath)
	_ = asciiArt
	fmt.Println("ASCII Art generation completed.")

}

func getImagePath(filename string) string {
	baseDir := os.Getenv("ASCII_NGIN_IMAGES_DIR")
	if baseDir == "" {
		baseDir = "images" // fallback to relative path
	}
	return filepath.Join(baseDir, filename)
}

func ExtractFrames() (VideoMetadata, error) {
	return VideoMetadata{}, nil
}

// in future we shouldnt need to pass in the filename we can get it based on the framenumber
func FrameToASCII(frame FrameData, filename string) [][]rune {
	height, err := GetImageHeight(filename)
	if err != nil {
		fmt.Println("Error getting image height:", err)
		return [][]rune{}
	}

	newHeight := ascii.RecalculateHeight(height)

	// Get just the base filename and create output paths in appropriate directories
	dir := filepath.Dir(filename)
	base := filepath.Base(filename)
	resizedFilename := filepath.Join(dir, "resized_images", "resized_"+base)
	resizeerr := ResizeImage(filename, resizedFilename, newHeight)
	if resizeerr != nil {
		fmt.Println("Error resizing image:", resizeerr)
		return [][]rune{}
	}

	textPixelDataFile := filepath.Join(dir, "text_frames", "text_"+base+".txt")
	err = ConvertToTextPixelData(resizedFilename, textPixelDataFile)
	if err != nil {
		fmt.Println("Error converting image to text pixel data:", err)
		return [][]rune{}
	}

	// Here we would read the textPixelDataFile and convert each pixel to ASCII
	// For simplicity, we'll return an empty 2D slice for now
	return [][]rune{}
}

// given rbg values of a pixel return the corresponding ascii character
func PixelToASCII(pixel PixelData) rune {

	luminance := ascii.CalculateLuminance(pixel.Red, pixel.Green, pixel.Blue)

	if luminance < ascii.LUMINANCE_THRESHOLD {
		return ' '
	}

	return ascii.MapLuminanceToASCII(luminance)
}

func GetFrameCount(video VideoData) int {
	return video.FrameCount
}
func GetFrameRate(video VideoData) float64 {
	return video.FrameRate
}
func GetFrames(video VideoData) []FrameData {
	return video.Frames
}
