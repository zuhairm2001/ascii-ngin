package video

import (
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/term"
)

// CheckDependencies verifies required system dependencies are installed
func CheckDependencies() error {
	err := checkFFmpegInstalled()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func GetTerminalDimensions() (TermSize, error) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return TermSize{}, err
	}
	return TermSize{Width: width, Height: height}, nil
}

func checkFFmpegInstalled() error {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("ffmpeg is not installed or not found in PATH")
	}
	return nil
}
