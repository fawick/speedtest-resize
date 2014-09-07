// +build !noexec

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func init() {
	// Program checks are supposed to work on all systems
	if _, err := exec.LookPath("gm"); err == nil {
		RegisterResizer("GraphicsMagick_thumbnail", graphicsMagickThumbnail)
	} else {
		fmt.Println("Cannot find gm in PATH, will skip GraphicsMagick tests")
	}
	if _, err := exec.LookPath("convert"); err == nil {
		RegisterResizer("ImageMagick_thumbnail", imageMagickThumbnail)
		RegisterResizer("ImageMagick_resize", imageMagickResize)
	} else {
		fmt.Println("Cannot find convert in PATH, will skip ImageMagick tests")
	}

}

func imageMagickThumbnail(origName, newName string) (int, int64) {
	origFileStat, _ := os.Stat(origName)
	var args = []string{
		"-define", "jpeg:size=300x300",
		"-thumbnail", "150x150>",
		origName, newName,
	}

	var cmd *exec.Cmd
	path, _ := exec.LookPath("convert")
	cmd = exec.Command(path, args...)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return 0, origFileStat.Size()
	}
	newFileStat, _ := os.Stat(newName)
	return int(newFileStat.Size()), origFileStat.Size()
}

func imageMagickResize(origName, newName string) (int, int64) {
	origFileStat, _ := os.Stat(origName)
	var args = []string{
		"-resize", "150x150>",
		origName, newName,
	}

	var cmd *exec.Cmd
	path, _ := exec.LookPath("convert")
	cmd = exec.Command(path, args...)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return 0, origFileStat.Size()
	}
	newFileStat, _ := os.Stat(newName)
	return int(newFileStat.Size()), origFileStat.Size()
}

func graphicsMagickThumbnail(origName, newName string) (int, int64) {
	origFileStat, _ := os.Stat(origName)
	var args = []string{
		"convert",
		"-define", "jpeg:size=300x300",
		"-thumbnail", "150x150>",
		origName, newName,
	}

	var cmd *exec.Cmd
	path, _ := exec.LookPath("gm")
	cmd = exec.Command(path, args...)

	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return 0, origFileStat.Size()
	}
	newFileStat, _ := os.Stat(newName)
	return int(newFileStat.Size()), origFileStat.Size()
}
