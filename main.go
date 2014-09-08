package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
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

type ResizerStat struct {
	Desc      string        // Description string
	Total     time.Duration // Total duration for all files
	Processed int           // Number of processed files
}

func (s ResizerStat) Avg() time.Duration {
	return s.Total / time.Duration(s.Processed)
}

type ResizerStats []*ResizerStat

func (s ResizerStats) Len() int      { return len(s) }
func (s ResizerStats) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByAverage struct{ ResizerStats }

func (b ByAverage) Less(i, j int) bool { return b.ResizerStats[i].Avg() < b.ResizerStats[j].Avg() }

func (s ResizerStats) WriteTo(w io.Writer) {
	formatHeader := "| %-30s |%-17s|%-9s|\n"
	formatRow := "| %-30s | %15.3fs | %-7s |\n"
	fmt.Fprintf(w, "\nResults\n-------\n\n")
	fmt.Fprintf(w, formatHeader, "Table", " Time (file avg.) ", " Pure Go ")
	fmt.Fprintf(w, formatHeader, strings.Repeat("-", 30), " ----------------:", ":-------:")
	for _, st := range s {
		fmt.Fprintf(w, formatRow, st.Desc, float64(st.Avg())/1e9, "")
	}
	fmt.Fprintln(w)
}

type ResizerFunc func(oldName, newName string) (int, int64)

type Resizer struct {
	Desc string // Description string
	// Func opens the image with the old name, resizes it, and saves it under
	// the new name. It returns the old and the new file size.
	Func ResizerFunc
}

func (r Resizer) Resize(files []string) *ResizerStat {
	s := ResizerStat{Desc: r.Desc}
	for i, origPath := range files {
		newPath := fmt.Sprintf("%s.thumb.%s.jpg", origPath, r.Desc)
		if *verbose {
			fmt.Printf("File %d w/ %s: ", i+1, s.Desc)
		}
		imgStart := time.Now()
		n, o := r.Func(origPath, newPath)
		ratio := float64(n) / float64(o) * 100.0
		dur := time.Since(imgStart)
		s.Processed++
		s.Total += dur
		avg := s.Total / time.Duration(s.Processed)
		if *verbose {
			fmt.Printf("re-encoded to size=%d (%.1f%%) in %s. New avg=%s\n", n, ratio, dur, avg)
		}
	}
	return &s
}

var RegisteredResizers []Resizer

func RegisterResizer(d string, f func(oldName, newName string) (int, int64)) {
	RegisteredResizers = append(RegisteredResizers, Resizer{Desc: d, Func: f})
}

var verbose = flag.Bool("verbose", false, "Print statistics for every single file processed")

func main() {
	flag.Parse()
	dir := "."
	if len(flag.Args()) > 1 {
		fmt.Println(flag.Args())
		dir = flag.Args()[1]
	}
	files, _ := scanDir(dir)
	if len(files) == 0 {
		fmt.Println("no jpg files found in", dir)
		return
	}
	if len(files) > 10 {
		files = files[0:10]
	}

	var results ResizerStats
	for _, r := range RegisteredResizers {
		results = append(results, r.Resize(files))
	}

	sort.Sort(ByAverage{results})
	results.WriteTo(os.Stdout)
}
