package main

import (
	"image"
	"image/color"
)

// Brightness is a filter that modifies the brightness of the image
type Brightness struct {
	Factor uint32 // Percentage of brightness to apply
}

// Process applies a brightness filter to the image
func (filter *Brightness) Process(in image.Image, out *image.RGBA, bounds image.Rectangle) {
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := in.At(x, y).RGBA()

  		        // r, g, b, a are 16bits components in a uint32
			nr := Trunc(r + (0xFFFF * filter.Factor) / 100)
			ng := Trunc(g + (0xFFFF * filter.Factor) / 100)
			nb := Trunc(b + (0xFFFF * filter.Factor) / 100)

			nc := color.NRGBA64{nr, ng, nb, uint16(a)}
			out.Set(x, y, nc)
		}
	}
}

// NewBrightness creates a new filter for brightness
func NewBrightness(factor uint32) *Brightness {
	return &Brightness{
		factor,
	}
}
