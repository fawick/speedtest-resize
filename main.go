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
	Resizer    Resizer
	Total      time.Duration // Total duration for all files
	Processed  int           // Number of processed files
	PercentSum float64
}

func (s ResizerStat) TimeAvg() time.Duration {
	return s.Total / time.Duration(s.Processed)
}

func (s ResizerStat) SizeAvg() float64 {
	return s.PercentSum / float64(s.Processed)
}

type ResizerStats []*ResizerStat

func (s ResizerStats) Len() int      { return len(s) }
func (s ResizerStats) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByAverage struct{ ResizerStats }

func (b ByAverage) Less(i, j int) bool {
	return b.ResizerStats[i].TimeAvg() < b.ResizerStats[j].TimeAvg()
}

func (s ResizerStats) WriteTo(w io.Writer) {
	formatHeader := "|%-32s|%-18s|%-18s|%-9s|\n"
	formatRow := "| %-30s | %15.3fs | %15.3f%% | %-7s |\n"
	fmt.Fprintf(w, "\nResults\n-------\n\n")
	fmt.Fprintf(w, formatHeader, " Table", " Time (file avg.)", " Size (file avg.)", " Pure Go")
	fmt.Fprintf(w, formatHeader, strings.Repeat("-", 32), "-----------------:", "-----------------:", ":-------:")
	for _, st := range s {
		var pureFlag = ""
		if st.Resizer.Pure {
			pureFlag = "   X"
		}
		fmt.Fprintf(w, formatRow, st.Resizer.Name, float64(st.TimeAvg())/1e9, st.SizeAvg(), pureFlag)
	}
	fmt.Fprintln(w)
}

type ResizerFunc func(oldName, newName string) (int, int64)

type Resizer struct {
	// Name or description
	Name string
	// Func opens the image with the old name, resizes it, and saves it under
	// the new name. It returns the old and the new file size.
	Func ResizerFunc
	// true iff Resizer has no dependency on non-Go code or external programs
	Pure bool
}

func (r Resizer) Resize(files []string) *ResizerStat {
	s := ResizerStat{Resizer: r}
	for i, origPath := range files {
		newPath := fmt.Sprintf("%s.thumb.%s.jpg", origPath, r.Name)
		if *verbose {
			fmt.Printf("File %d w/ %s: ", i+1, r.Name)
		}
		imgStart := time.Now()
		n, o := r.Func(origPath, newPath)
		ratio := float64(n) / float64(o) * 100.0
		dur := time.Since(imgStart)
		s.PercentSum += ratio
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

func RegisterResizer(n string, f func(oldName, newName string) (int, int64)) {
	RegisteredResizers = append(RegisteredResizers, Resizer{Name: n, Func: f})
}

func RegisterPureResizer(n string, f func(oldName, newName string) (int, int64)) {
	RegisteredResizers = append(RegisteredResizers, Resizer{n, f, true})
}

var verbose = flag.Bool("verbose", false, "Print statistics for every single file processed")

func main() {
	flag.Parse()
	dir := "."
	if len(flag.Args()) > 0 {
		dir = flag.Args()[0]
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
