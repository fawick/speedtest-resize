// +build !nopure

package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"github.com/bamiaux/rez"
	"github.com/disintegration/gift"
	"github.com/disintegration/imaging"
	nfnt_resize "github.com/nfnt/resize"

	"golang.org/x/image/draw"
)

func init() {
	RegisterPureResizer("imaging_box", resizeImaging)
	RegisterPureResizer("gift_box", resizeGift)
	RegisterPureResizer("Nfnt_NearestNeighbor", resizeNfntNearestNeighbor)
	RegisterPureResizer("rez_bilinear", resizeRezBilinear)
	RegisterPureResizer("x_image_draw", resizeXImageDraw)
}

func resizeNfnt(origName, newName string, interp nfnt_resize.InterpolationFunction) (int, int64) {
	origFile, _ := os.Open(origName)
	origImage, _ := jpeg.Decode(origFile)
	origFileStat, _ := origFile.Stat()
	origFile.Close()

	var resized image.Image
	resized = nfnt_resize.Thumbnail(150, 150, origImage, interp)

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

func resizeNfntNearestNeighbor(origName, newName string) (int, int64) {
	return resizeNfnt(origName, newName, nfnt_resize.NearestNeighbor)
}

func getSize(a, b, c int) int {
	d := a * b / c
	return (d + 1) & -1
}

func resizeRez(origName, newName string, filter rez.Filter) (int, int64) {
	origFile, _ := os.Open(origName)
	origImage, _ := jpeg.Decode(origFile)
	origFileStat, _ := origFile.Stat()
	origFile.Close()

	var resized image.Image
	src, ok := origImage.(*image.YCbCr)
	if !ok {
		fmt.Println("input picture is not ycbcr")
		return 0, origFileStat.Size()
	}

	p := origImage.Bounds().Size()
	w, h := 150, getSize(150, p.Y, p.X)
	if p.X < p.Y {
		w, h = getSize(150, p.X, p.Y), 150
	}
	resized = image.NewYCbCr(image.Rect(0, 0, w, h), src.SubsampleRatio)
	err := rez.Convert(resized, origImage, filter)
	if err != nil {
		fmt.Println("unable to convert picture", err)
		return 0, origFileStat.Size()
	}

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

func resizeRezBilinear(origName, newName string) (int, int64) {
	return resizeRez(origName, newName, rez.NewBilinearFilter())
}

func resizeImaging(origName, newName string) (int, int64) {
	origFileStat, _ := os.Stat(origName)
	origImage, _ := imaging.Open(origName)

	var resized image.Image
	resized = imaging.Fit(origImage, 150, 150, imaging.Box)

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

func resizeGift(origName, newName string) (int, int64) {
	origFile, _ := os.Open(origName)
	origImage, _ := jpeg.Decode(origFile)
	origFileStat, _ := origFile.Stat()
	origFile.Close()

	var g = gift.New()

	p := origImage.Bounds().Size()
	if p.X > p.Y {
		g.Add(gift.Resize(150, 0, gift.BoxResampling))
	} else {
		g.Add(gift.Resize(0, 150, gift.BoxResampling))
	}
	resized := image.NewRGBA(g.Bounds(origImage.Bounds()))
	g.Draw(resized, origImage)

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

func resizeXImageDraw(origName, newName string) (int, int64) {
	origFile, _ := os.Open(origName)
	origImage, _ := jpeg.Decode(origFile)
	origFileStat, _ := origFile.Stat()
	origFile.Close()

	p := origImage.Bounds().Size()
	w, h := 150, getSize(150, p.Y, p.X)
	if p.X < p.Y {
		w, h = getSize(150, p.X, p.Y), 150
	}
	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	draw.Draw(dst, dst.Bounds(), image.White, image.ZP, draw.Src)
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), origImage, origImage.Bounds(), draw.Src, nil)

	b := new(bytes.Buffer)
	jpeg.Encode(b, dst, nil)
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
