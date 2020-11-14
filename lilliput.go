// +build lilliput all

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/discordapp/lilliput"
)

func init() {
	RegisterResizer("lilliput", resizeLilliput)
}

func resizeLilliPut(origName, newName string) (int, int64) {
	origFileStat, _ := os.Stat(origName)
	buf, err := ioutil.ReadFile(origName)
	if err != nil {
		fmt.Println(err)
		return 0, origFileStat.Size()
	}
	d, err := lilliput.NewDecoder(buf)
	if err != nil {
		fmt.Println(err)
		return 0, origFileStat.Size()
	}
	defer d.Close()

	ops := lilliput.NewImageOps(8192)
	defer ops.Close()

	// create a buffer to store the output image, 50MB in this case
	outputImg := make([]byte, 50*1024*1024)

	resizeMethod := lilliput.ImageOpsFit
	if stretch {
		resizeMethod = lilliput.ImageOpsResize
	}

	if outputWidth == header.Width() && outputHeight == header.Height() {
		resizeMethod = lilliput.ImageOpsNoResize
	}

	opts := &lilliput.ImageOptions{
		FileType:             filepath.Ext(newName),
		Width:                150,
		Height:               150,
		ResizeMethod:         lilliput.ImageOpsFit,
		NormalizeOrientation: true,
		EncodeOptions:        map[int]int{lilliput.JpegQuality: 95},
	}

	// resize and transcode image
	outputImg, err = ops.Transform(decoder, opts, outputImg)
	if err != nil {
		fmt.Printf("error transforming image, %s\n", err)
		return 0, origFileStat.Size()
	}

	// image has been resized, now write file out
	if outputFilename == "" {
		return 0, origFileStat.Size()
	}

	err = ioutil.WriteFile(outputFilename, outputImg, 0600)
	if err != nil {
		fmt.Printf("error writing out resized image, %s\n", err)
		return 0, origFileStat.Size()
	}
	return len(outputImg), origFileStat.Size()

}
