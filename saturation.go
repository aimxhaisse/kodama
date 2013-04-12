package main

import (
	"image"
	"image/color"
)

// Saturation is a filter that modifies the saturation of the image
type Saturation struct {
	Factor uint32 // Percentage of saturation to apply
}

// Process applies a darkness filter to the image
func (filter *Saturation) Process(in image.Image, out *image.RGBA, bounds image.Rectangle) {
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := in.At(x, y).RGBA()

			grey := (r + g + b) / 3

			nr := Trunc(r + ((Abs(int32(r-grey)) * filter.Factor) / 100))
			ng := Trunc(g + ((Abs(int32(g-grey)) * filter.Factor) / 100))
			nb := Trunc(b + ((Abs(int32(b-grey)) * filter.Factor) / 100))

			nc := color.NRGBA64{nr, ng, nb, uint16(a)}
			out.Set(x, y, nc)
		}
	}
}

// NewSaturation creates a new filter for darkness
func NewSaturation(factor uint32) *Saturation {
	return &Saturation{
		factor,
	}
}
