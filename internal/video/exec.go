package video

import (
	"fmt"
	"os/exec"
)

func GetImageHeight(filename string) (int, error) {
	out, err := exec.Command("magick", "identify", "-format", "%h", filename).Output()
	if err != nil {
		return 0, err
	}

	var height int
	_, err = fmt.Sscanf(string(out), "%d", &height)
	if err != nil {
		return 0, err
	}

	return height, nil
}

func ResizeImage(inputFile string, outputFile string, newHeight int) error {
	stdout, err := exec.Command("magick", inputFile, "-resize", fmt.Sprintf("x%d", newHeight), outputFile).Output()
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
