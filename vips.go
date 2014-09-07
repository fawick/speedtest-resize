// +build vips all

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/DAddYE/vips"
)

func init() {
	RegisterResizer("vips", resizeVips)
}

func resizeVips(origName, newName string) (int, int64) {
	origFileStat, _ := os.Stat(origName)
	options := vips.Options{
		Width:        150,
		Height:       150,
		Crop:         false,
		Extend:       vips.EXTEND_WHITE,
		Interpolator: vips.BILINEAR,
		Gravity:      vips.CENTRE,
		Quality:      95,
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
	buf, err = vips.Resize(buf, options)
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
