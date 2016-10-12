package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"../common"

	_ "image/jpeg"
)

import _ "image/gif"

func main() {
	var arg string
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}
	switch arg {
	case "", "-h", "--help":
		fmt.Println("Usage: fetcher <PATH_TO_PHOTO_FOLDER>")
	default:
		start := time.Now()
		colorsFile := filepath.Join(arg, "colors.txt")
		if files, err := fetchFileList(arg); err == nil {
			if list, err := fetchColors(files); err == nil {
				data := []byte(strings.Join(list, "\n"))
				if err := ioutil.WriteFile(colorsFile, data, 0644); err != nil {
					fmt.Println("Couldn't write file:", err)
				}
				fmt.Printf("\nPicture colors fetched and dumped to %s (%v)\n", colorsFile, time.Since(start))
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
		if err != nil {
			return err
		}
		// ignore non-jpg's
		if !strings.HasSuffix(info.Name(), ".jpg") {
			return nil
		}
		list = append(list, path)
		return nil
	})
	return list, err
}

func fetchColors(files []string) (list []string, err error) {
	for _, path := range files {
		r, g, b, a, err := averageColor(path)
		if err != nil {
			return nil, fmt.Errorf("couldn't extract colors from '%v': %v", path, err)
		}
		list = append(list, path, fmt.Sprintf("%v %v %v %v", r, g, b, a))
	}
	return list, err
}

func averageColor(path string) (uint, uint, uint, uint, error) {
	m, err := common.ReadImage(path)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return common.AverageColorFromBounds(m, m.Bounds())
}
