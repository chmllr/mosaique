package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
		r, g, b, a, err := getAVGColor(path)
		if err != nil {
			return fmt.Errorf("couldn't extract colors from '%v': %v", path, err)
		}
		list = append(list, path, fmt.Sprintf("%v %v %v %v", r, g, b, a))
		return nil
	})
	return list, err
}

func getAVGColor(path string) (uint, uint, uint, uint, error) {
	reader, err := os.Open(path)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	defer reader.Close()
	m, _, err := image.Decode(reader)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("decoding failed: %v", err)
	}
	bounds := m.Bounds()
	size := uint64((bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y))
	fmt.Println("Reading", path, "with size", size)
	var r, g, b, a uint64
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r32, g32, b32, a32 := m.At(x, y).RGBA()
			r += uint64(r32)
			g += uint64(g32)
			b += uint64(b32)
			a += uint64(a32)
		}
	}
	return uint(r / size), uint(g / size), uint(b / size), uint(a / size), nil
}
