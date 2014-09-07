package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"strings"
	"time"
)

func scanDir(path string) (files []string, hello error) {
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

type ResizerFunc func(oldName, newName string) (int, int64)

type Resizer struct {
	Desc string // Description string
	// Func opens the image with the old name, resizes it, and saves it under
	// the new name. It returns the old and the new file size.
	Func ResizerFunc
}

func (r Resizer) Resize(files []string) string {
	start := time.Now()

	var total int64

	for i, origPath := range files {
		newPath := fmt.Sprintf("%s.thumb.%s.jpg", origPath, r.Desc)
		if printSingleFile {
			fmt.Printf("File %d: ", i)
		}
		imgStart := time.Now()
		n, o := r.Func(origPath, newPath)
		ratio := float64(n) / float64(o) * 100.0
		dur := time.Since(imgStart)
		total += int64(dur)
		avg := time.Duration(total / int64(i+1))
		if printSingleFile {
			fmt.Printf("re-encoded to size=%d (%.1f%%) in %s. New avg=%s\n", n, ratio, dur, avg)
		}
	}
	return timeTrack(start, r.Desc, len(files))
}

var printSingleFile bool

var RegisteredResizers []Resizer

func RegisterResizer(d string, f func(oldName, newName string) (int, int64)) {
	RegisteredResizers = append(RegisteredResizers, Resizer{Desc: d, Func: f})
}

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

	for _, r := range RegisteredResizers {
		results = append(results, r.Resize(files))
	}

	for _, s := range results {
		fmt.Println(s)
	}
}
