// +build !noexec

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func init() {
	// Program checks are supposed to work on all systems
	if _, err := exec.LookPath("gm"); err == nil {
		RegisterResizer("GraphicsMagick_thumbnail", graphicsMagickThumbnail)
	}
	if _, err := exec.LookPath("convert"); err == nil {
		RegisterResizer("ImageMagick_thumbnail", imageMagickThumbnail)
		RegisterResizer("ImageMagick_resize", imageMagickResize)
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
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("/usr/bin/convert", args...)
	case "windows":
		path, _ := exec.LookPath("convert.exe")
		cmd = exec.Command(path, args...)
	}
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
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("/usr/bin/convert", args...)
	case "windows":
		path, _ := exec.LookPath("convert.exe")
		cmd = exec.Command(path, args...)
	}
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
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("/usr/bin/gm", args...)
	case "windows":
		path, _ := exec.LookPath("gm.exe")
		cmd = exec.Command(path, args...)
	}

	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return 0, origFileStat.Size()
	}
	newFileStat, _ := os.Stat(newName)
	return int(newFileStat.Size()), origFileStat.Size()
}
