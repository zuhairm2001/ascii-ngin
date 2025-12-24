package video

import (
	"fmt"
	"os/exec"
)

func GetImageDimensions(filename string) (int, int, error) {
	out, err := exec.Command("magick", "identify", "-format", "%hx%w", filename).Output()
	if err != nil {
		return 0, 0, err
	}

	var height, width int
	_, err = fmt.Sscanf(string(out), "%dx%d", &height, &width)
	if err != nil {
		return 0, 0, err
	}

	return height, width, nil
}

func ResizeImage(inputFile string, outputFile string, newWidth int, newHeight int) error {
	stdout, err := exec.Command("magick", inputFile, "-resize", fmt.Sprintf("%dx%d", newWidth, newHeight), outputFile).Output()
	fmt.Print(string(stdout))
	fmt.Print(err)
	if err != nil {
		return err
	}
	return nil
}

func ConvertToTextPixelData(imageFile string, outputFile string) error {
	_, err := exec.Command("magick", imageFile, "txt:"+outputFile).Output()
	if err != nil {
		return err
	}
	return nil
}
