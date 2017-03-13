// +build imagick all

package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/gographics/imagick.v2/imagick"
)

func init() {
	imagick.Initialize()
	// TODO send imagick.Terminate to main() for clean deconstruction

	RegisterResizer("magickwand_box", resizeMagickWand)
}

func resizeMagickWand(origName, newName string) (int, int64) {
	origFileStat, _ := os.Stat(origName)
	var err error
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err = mw.ReadImage(origName)
	if err != nil {
		fmt.Println(err)
		return 0, origFileStat.Size()
	}
	start := time.Now()

	filter := imagick.FILTER_BOX
	w := mw.GetImageWidth()
	h := mw.GetImageHeight()
	if w > h {
		err = mw.ResizeImage(150, 100, filter, 1)
	} else {
		err = mw.ResizeImage(100, 150, filter, 1)
	}
	if err != nil {
		fmt.Println(time.Since(start))
		fmt.Println(err)
		return 0, origFileStat.Size()
	}

	err = mw.SetImageCompressionQuality(95)
	if err != nil {
		fmt.Println(err)
		return 0, origFileStat.Size()
	}

	err = mw.WriteImage(newName)
	if err != nil {
		fmt.Println(err)
		return 0, origFileStat.Size()
	}

	newFileStat, _ := os.Stat(newName)
	return int(newFileStat.Size()), origFileStat.Size()
}
