// +build opencv all

package main

import (
	"os"

	"github.com/lazywei/go-opencv/opencv"
)

func init() {
	RegisterResizer("opencv", resizeOpenCv)
}

func resizeOpenCv(origName, newName string) (int, int64) {
	iplImg := opencv.LoadImage(origName)
	if iplImg == nil {
		panic("LoadImage fail")
	}
	defer iplImg.Release()
	resizedIplImg := opencv.Resize(iplImg, 150, 0, 0)
	opencv.SaveImage(newName, resizedIplImg, 0)

	origFileStat, _ := os.Stat(origName)
	newFileStat, _ := os.Stat(newName)

	return int(newFileStat.Size()), origFileStat.Size()
}
