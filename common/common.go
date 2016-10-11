package common

import "image"

// AverageColorFromBounds returns the average color for the given rectangle
func AverageColorFromBounds(m image.Image, bounds image.Rectangle) (uint, uint, uint, uint, error) {
	size := uint64((bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y))
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
