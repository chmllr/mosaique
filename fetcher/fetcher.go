package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"../common"
)

func main() {
	var arg string
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}
	switch arg {
	case "", "-h", "--help":
		fmt.Println("Usage: fetcher <PATH_TO_PHOTO_FOLDER>")
	default:
		fmt.Println("Scanning path", arg)
		start := time.Now()
		colorsFile := "colors.txt"
		if files, err := fetchFileList(arg); err == nil {
			fmt.Println(len(files), "files found")
			if list, err := fetchColors(files); err == nil {
				data := []byte(strings.Join(list, "\n"))
				if err := ioutil.WriteFile(colorsFile, data, 0644); err != nil {
					fmt.Println("Couldn't write file:", err)
				}
				fmt.Printf("Picture colors fetched and dumped to %s (%v)\n", colorsFile, time.Since(start))
			} else {
				fmt.Println("Couldn't fetch colors:", err)
			}
		} else {
			fmt.Println("Couldn't fetch file list:", err)
		}
	}
}

func fetchFileList(path string) (list []string, err error) {
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(strings.ToLower(path), ".jpg") {
			list = append(list, path)
		}
		return err
	})
	return list, err
}

func fetchColors(files []string) (list []string, err error) {
	results := make(chan *common.Entry)
	cpus := runtime.NumCPU()
	fmt.Println("Starting", cpus, "concurrent routines")
	for cpu := 0; cpu < cpus; cpu++ {
		go func(index int) {
			fmt.Println("Routine", index, "started...")
			for i := index; i < len(files); i += cpus {
				if entry, err := averageColor(files[i]); err == nil {
					results <- entry
				} else {
					fmt.Println("Error: couldn't get colors from", files[i])
				}
			}
		}(cpu)
	}
	for range files {
		c := <-results
		list = append(list, c.Path, fmt.Sprintf("%v %v %v %v", c.R, c.G, c.B, c.A))
	}
	close(results)
	return list, err
}

func averageColor(path string) (*common.Entry, error) {
	m, err := common.ReadImage(path)
	if err != nil {
		return nil, err
	}
	r, g, b, a, err := common.AverageColorFromBounds(m, m.Bounds())
	return &common.Entry{R: r, G: g, B: b, A: a, Path: path}, err
}
