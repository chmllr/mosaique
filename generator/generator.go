package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"../common"
)

const (
	tileSize = 28 // px
)

func printHelp() {
	fmt.Println("Usage: generator <PATH_TO_COLOR_FILE> <PATH_TO_SOURCE_IMAGE>")
}

func main() {
	var data, source string
	if len(os.Args) > 2 {
		data = os.Args[1]
		source = os.Args[2]
	} else {
		printHelp()
		return
	}
	switch data {
	case "", "-h", "--help":
		printHelp()
	default:
		start := time.Now()
		srcImg, err := common.ReadImage(source)
		if err != nil {
			fmt.Println("Couldn't read source image:", err)
			return
		}
		colors, err := makeColorVectors(data)
		if err != nil {
			fmt.Println("Couldn't create color vectors:", err)
			return
		}
		mos, err := createMosaique(srcImg, colors)
		if err != nil {
			fmt.Println("Couldn't create mosaique:", err)
			return
		}
		mosPath := source + "_mosaique.jpg"
		mosFile, err := os.Create(mosPath)
		jpeg.Encode(mosFile, mos, &jpeg.Options{Quality: jpeg.DefaultQuality})
		fmt.Printf("Mosaique written to %s (%s)\n", mosPath, time.Since(start))
	}
}

func createMosaique(orig image.Image, colors []*common.Entry) (image.Image, error) {
	bounds := orig.Bounds()
	mos := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X-tileSize; x += tileSize {
		for y := bounds.Min.Y; y < bounds.Max.Y-tileSize; y += tileSize {
			tileHolder := image.Rect(x, y, x+tileSize, y+tileSize)
			r, g, b, a, err := common.AverageColorFromBounds(orig, tileHolder)
			if err != nil {
				return nil, err
			}
			closestTile := findClosestTile(colors, r, g, b, a)
			tile, err := common.ReadImage(closestTile.Path)
			if err != nil {
				return nil, err
			}
			draw.Draw(mos, tileHolder, tile, image.ZP, draw.Src)
		}
	}
	return mos, nil
}

func findClosestTile(colors []*common.Entry, r, g, b, a uint16) *common.Entry {
	var res *common.Entry
	minDist := uint64(math.MaxUint64)
	for _, e := range colors {
		dist := getDist(r, g, b, a, e)
		if dist < minDist {
			minDist = dist
			res = e
		}
	}
	if res == nil {
		panic("unexpected situation: no closest tile found")
	}
	return res
}

func getDist(r, g, b, a uint16, e *common.Entry) uint64 {
	return p2(r-e.R) + p2(g-e.G) + p2(b-e.B) + p2(a-e.A)
}

func p2(a uint16) uint64 { return uint64(a) * uint64(a) }

func makeColorVectors(path string) ([]*common.Entry, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	res := []*common.Entry{}
	lines := strings.Split(string(data), "\n")
	for i := 0; i < len(lines)-1; i += 2 {
		colors := strings.Split(lines[i+1], " ")
		r, _ := strconv.Atoi(colors[0])
		g, _ := strconv.Atoi(colors[1])
		b, _ := strconv.Atoi(colors[2])
		a, _ := strconv.Atoi(colors[3])
		res = append(res, &common.Entry{Path: lines[i], R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a)})
	}
	return res, nil
}