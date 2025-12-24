package video

import (
	"fmt"
	"os/exec"
)

// making sure all of our deps are installed, assuming UNIX like system
func CheckDependencies() error {
	err := checkFFmpegInstalled()
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = checkMagickInstalled()
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = checkAwkInstalled()
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = checkFindInstalled()
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = checkTailInstalled()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

func checkFFmpegInstalled() error {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("ffmpeg is not installed or not found in PATH")
	}
	return nil
}

func checkMagickInstalled() error {
	_, err := exec.LookPath("magick")
	if err != nil {
		return fmt.Errorf("ImageMagick is not installed or not found in PATH")
	}
	return nil
}

func checkAwkInstalled() error {
	_, err := exec.LookPath("awk")
	if err != nil {
		return fmt.Errorf("awk is not installed or not found in PATH")
	}
	return nil
}

func checkFindInstalled() error {
	_, err := exec.LookPath("find")
	if err != nil {
		return fmt.Errorf("find is not installed or not found in PATH")
	}
	return nil
}

func checkTailInstalled() error {
	_, err := exec.LookPath("tail")
	if err != nil {
		return fmt.Errorf("tail is not installed or not found in PATH")
	}
	return nil
}
