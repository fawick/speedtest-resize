// +build opencv all

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/daddye/trez"
	"github.com/lazywei/go-opencv/opencv"
)

func init() {
	RegisterResizer("opencv", resizeOpenCv)
	RegisterResizer("T-REZ", resizeTRez)
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

func resizeTRez(origName, newName string) (int, int64) {
	origFileStat, _ := os.Stat(origName)
	options := trez.Options{
		Width:   150,
		Height:  150,
		Algo:    trez.FIT,
		Quality: 95,
	}
	origFile, err := os.Open(origName)
	if err != nil {
		fmt.Println(err)
		return 0, origFileStat.Size()
	}
	defer origFile.Close()
	buf, err := ioutil.ReadAll(origFile)
	if err != nil {
		fmt.Println(err)
		return 0, origFileStat.Size()
	}
	buf, err = trez.Resize(buf, options)
	if err != nil {
		fmt.Println(err)
		return 0, origFileStat.Size()
	}
	cacheFile, err := os.Create(newName)
	if err != nil {
		fmt.Println(err)
		return 0, origFileStat.Size()
	}
	defer cacheFile.Close()
	_, err = cacheFile.Write(buf)
	if err != nil {
		fmt.Println(err)
		return 0, origFileStat.Size()
	}

	return int(len(buf)), origFileStat.Size()
}
