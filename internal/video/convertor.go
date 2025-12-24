package video

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zuhairm2001/ascii-ngin/pkg/ascii"
)

type PixelData struct {
	Red   int
	Green int
	Blue  int
	X     int
	Y     int
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
	PrintASCIIArt(asciiArt)
	err = PrintASCIIArtToFile(asciiArt, "output.txt")
	if err != nil {
		fmt.Println("Error writing ASCII art to file:", err)
		return
	}
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
	// TODO
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

	pixelData, err := ReadTextFile(textPixelDataFile)
	if err != nil {
		fmt.Println("Error reading text pixel data file:", err)
		return [][]rune{}
	}

	var asciiArt [][]rune
	for _, row := range pixelData {
		var asciiRow []rune
		for _, pixel := range row {
			asciiChar := PixelToASCII(pixel)
			asciiRow = append(asciiRow, asciiChar)
		}
		asciiArt = append(asciiArt, asciiRow)
	}

	return asciiArt
}

// given rbg values of a pixel return the corresponding ascii character
func PixelToASCII(pixel PixelData) rune {

	luminance := ascii.CalculateLuminance(pixel.Red, pixel.Green, pixel.Blue)

	if luminance < ascii.LUMINANCE_THRESHOLD {
		return ' '
	}

	return ascii.MapLuminanceToASCII(luminance)
}

func ReadTextFile(filePath string) ([][]PixelData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var pixelData [][]PixelData
	var currentRowPixelArr []PixelData
	currentRow := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] == '#' {
			continue
		}

		pixel, err := ReadLine(line)
		if err != nil {
			continue
		}

		if pixel.Y != currentRow {
			if len(currentRowPixelArr) > 0 {
				pixelData = append(pixelData, currentRowPixelArr)
			}
			currentRowPixelArr = []PixelData{}
			currentRow = pixel.Y
		}

		currentRowPixelArr = append(currentRowPixelArr, pixel)
	}

	if len(currentRowPixelArr) > 0 {
		pixelData = append(pixelData, currentRowPixelArr)
	}

	return pixelData, nil
}

func ReadLine(line string) (PixelData, error) {
	var x, y, r, g, b int
	_, err := fmt.Sscanf(line, "%d,%d: (%d,%d,%d)", &x, &y, &r, &g, &b)
	if err != nil {
		return PixelData{}, err
	}
	return PixelData{X: x, Y: y, Red: r, Green: g, Blue: b}, nil
}

func PrintASCIIArt(asciiArt [][]rune) {
	for _, row := range asciiArt {
		for _, char := range row {
			fmt.Print(string(char))
		}
		fmt.Println()
	}
}

func PrintASCIIArtToFile(asciiArt [][]rune, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, row := range asciiArt {
		for _, char := range row {
			_, err := writer.WriteString(string(char))
			if err != nil {
				return err
			}
		}
		_, err := writer.WriteString("\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
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
