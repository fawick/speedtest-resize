// +build fastjpeg all

package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"os"

	"camlistore.org/pkg/images/fastjpeg"
	"camlistore.org/pkg/images/resize"
)

func init() {
	if fastjpeg.Available() {
		RegisterResizer("fastjpeg", resizeFastjpeg)
	} else {
		fmt.Println("Cannot find djpeg in PATH, will skip fastjpeg tests")
	}
}

func resizeFastjpeg(origName, newName string) (int, int64) {
	origFile, _ := os.Open(origName)
	origFileStat, _ := origFile.Stat()

	downsampled, err := fastjpeg.DecodeDownsample(origFile, 8)
	origFile.Close()
	if err != nil {
		fmt.Println(err)
		return 0, origFileStat.Size()
	}

	resized := resize.Resize(downsampled, downsampled.Bounds(), 150, 150)

	b := new(bytes.Buffer)
	jpeg.Encode(b, resized, nil)
	blen := b.Len()
	cacheFile, err := os.Create(newName)
	defer cacheFile.Close()
	if err != nil {
		fmt.Println(err)
		return 0, origFileStat.Size()
	}
	b.WriteTo(cacheFile)

	return blen, origFileStat.Size()
}
