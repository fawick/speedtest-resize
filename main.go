package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	moustaschio_resize "code.google.com/p/appengine-go/example/moustachio/resize"
	"github.com/bamiaux/rez"
	"github.com/disintegration/gift"
	"github.com/disintegration/imaging"
	"github.com/gographics/imagick/imagick"
	"github.com/lazywei/go-opencv/opencv"
	nfnt_resize "github.com/nfnt/resize"
)

func scanDir(path string) (files []string, err error) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, r := range entries {
		n := strings.ToUpper(r.Name())
		if strings.HasSuffix(n, ".JPG") || strings.HasSuffix(n, ".JPEG") {
			files = append(files, path+"/"+r.Name())
		}
	}
	return
}

func timeTrack(start time.Time, name string, n int) string {
	elapsed := time.Since(start)
	avg := time.Duration(int64(elapsed) / int64(n))
	s := fmt.Sprintf("%s took %s, file average %s\n", name, elapsed, avg)
	fmt.Println(s)
	return s
}

func resizeNfnt(origName, newName string, interp nfnt_resize.InterpolationFunction) (int, int64) {
	origFile, _ := os.Open(origName)
	origImage, _ := jpeg.Decode(origFile)
	origFileStat, _ := origFile.Stat()
	origFile.Close()

	var resized image.Image
	p := origImage.Bounds().Size()
	if p.X > p.Y {
		resized = nfnt_resize.Resize(150, 0, origImage, interp)
	} else {
		resized = nfnt_resize.Resize(0, 150, origImage, interp)
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

func moustachioResample(origName, newName string) (int, int64) {
	return resizeMoustachio(origName, newName, moustaschio_resize.Resample)
}

func moustachioResize(origName, newName string) (int, int64) {
	return resizeMoustachio(origName, newName, moustaschio_resize.Resize)
}

func resizeMoustachio(origName, newName string, method func(image.Image, image.Rectangle, int, int) image.Image) (int, int64) {
	origFile, _ := os.Open(origName)
	origImage, _ := jpeg.Decode(origFile)
	origFileStat, _ := origFile.Stat()
	origFile.Close()

	var resized image.Image
	p := origImage.Bounds().Size()
	if p.X > p.Y {
		resized = method(origImage, origImage.Bounds(), 150, 100)
	} else {
		resized = method(origImage, origImage.Bounds(), 100, 150)
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
	fmt.Print("resizing part: ", time.Since(start), ", ")

	newFileStat, _ := os.Stat(newName)
	return int(newFileStat.Size()), origFileStat.Size()
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

func resize(files []string, desc string, m func(string, string) (int, int64)) string {
	start := time.Now()

	var total int64

	for i, origPath := range files {
		newPath := fmt.Sprintf("%s.thumb.%s.jpg", origPath, desc)
		if printSingleFile {
			fmt.Printf("File %d: ", i)
		}
		imgStart := time.Now()
		n, o := m(origPath, newPath)
		ratio := float64(n) / float64(o) * 100.0
		dur := time.Since(imgStart)
		total += int64(dur)
		avg := time.Duration(total / int64(i+1))
		if printSingleFile {
			fmt.Printf("re-encoded to size=%d (%.1f%%) in %s. New avg=%s\n", n, ratio, dur, avg)
		}
	}
	return timeTrack(start, desc, len(files))
}

var printSingleFile bool

func main() {
	dir := "."
	if len(os.Args) > 1 {
		fmt.Println(os.Args)
		dir = os.Args[1]
	}
	printSingleFile = true
	files, _ := scanDir(dir)
	if len(files) == 0 {
		fmt.Println("no jpg files found in", dir)
		return
	}
	if len(files) > 10 {
		files = files[0:10]
	}
	var results []string

	imagick.Initialize()
	defer imagick.Terminate()
	results = append(results, resize(files, "magickwand_box", resizeMagickWand))

	switch runtime.GOOS {
	case "linux":
		if _, err := os.Stat("/usr/bin/gm"); err == nil {
			results = append(results, resize(files, "GraphicsMagick_thumbnail", graphicsMagickThumbnail))
		}

		if _, err := os.Stat("/usr/bin/convert"); err == nil {
			results = append(results, resize(files, "ImageMagick_thumbnail", imageMagickThumbnail))
			results = append(results, resize(files, "ImageMagick_resize", imageMagickResize))
		}
	case "windows":
		if _, err := exec.LookPath("gm.exe"); err == nil {
			results = append(results, resize(files, "GraphicsMagick_thumbnail", graphicsMagickThumbnail))
		}

		if _, err := exec.LookPath("convert.exe"); err == nil {
			results = append(results, resize(files, "ImageMagick_thumbnail", imageMagickThumbnail))
			results = append(results, resize(files, "ImageMagick_resize", imageMagickResize))
		}
	}

	results = append(results, resize(files, "imaging_Box", resizeImaging))
	results = append(results, resize(files, "gift_Box", resizeGift))
	results = append(results, resize(files, "moustaschio_resize", moustachioResize))
	results = append(results, resize(files, "Nfnt_NearestNeighbor", resizeNfntNearestNeighbor))
	results = append(results, resize(files, "OpenCv", resizeOpenCv))
	results = append(results, resize(files, "rez_bilinear", resizeRezBilinear))

	for _, s := range results {
		fmt.Println(s)
	}
}
