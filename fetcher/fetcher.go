package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"../common"
)

import _ "image/jpeg"
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
		if list, err := fetch(arg); err == nil {
			data := []byte(strings.Join(list, "\n"))
			if err := ioutil.WriteFile("colors.txt", data, 0644); err != nil {
				fmt.Println("Couldn't write file:", err)
			}
			fmt.Println("Picture colors fetched and dumped to ./colors.txt")
		} else {
			fmt.Println("Couldn't fetch colors:", err)
		}
	}
}

func fetch(path string) (list []string, err error) {
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// ignore non-jpg's
		if !strings.HasSuffix(info.Name(), ".jpg") {
			return nil
		}
		r, g, b, a, err := averageColor(path)
		if err != nil {
			return fmt.Errorf("couldn't extract colors from '%v': %v", path, err)
		}
		list = append(list, path, fmt.Sprintf("%v %v %v %v", r, g, b, a))
		return nil
	})
	return list, err
}

func averageColor(path string) (uint, uint, uint, uint, error) {
	m, err := common.ReadImage(path)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return common.AverageColorFromBounds(m, m.Bounds())
}
