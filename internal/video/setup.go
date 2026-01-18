package video

import (
	"fmt"
	"os/exec"
)

type TermData struct {
	Width  int
	Height int
}

// CheckDependencies verifies required system dependencies are installed
func CheckDependencies() error {
	err := checkFFmpegInstalled()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func GetTerminalDimensions() (TermData, error) {
	heightOut, err := exec.Command("tput", "lines").Output()
	widthOut, err := exec.Command("tput", "cols").Output()
	if err != nil {
		return TermData{}, err
	}

	var rows, cols int
	_, err = fmt.Sscanf(string(heightOut), "%d", &rows)
	if err != nil {
		return TermData{}, err
	}
	_, err = fmt.Sscanf(string(widthOut), "%d", &cols)
	if err != nil {
		return TermData{}, err
	}
	return TermData{Width: cols, Height: rows}, nil
}

func checkFFmpegInstalled() error {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("ffmpeg is not installed or not found in PATH")
	}
	return nil
}
