package filters

import (
	"image"
	"image/color"
)

// Darkness is a filter that modifies the darkness of the image
type Darkness struct {
	Factor uint32 // Percentage of darkness to apply
}

// IsScalable returns true because Darkness is a scalable filter
func (filter *Darkness) IsScalable() bool {
	return true
}

// Process applies a darkness filter to the image
func (filter *Darkness) Process(in image.Image, out *image.RGBA, bounds image.Rectangle) {
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := in.At(x, y).RGBA()

			nr := Strunc(int32(r - (0xFFFF*filter.Factor)/100))
			ng := Strunc(int32(g - (0xFFFF*filter.Factor)/100))
			nb := Strunc(int32(b - (0xFFFF*filter.Factor)/100))

			nc := color.NRGBA64{nr, ng, nb, uint16(a)}
			out.Set(x, y, nc)
		}
	}
}

// NewDarkness creates a new filter for darkness
func NewDarkness(factor uint32) *Darkness {
	return &Darkness{
		factor,
	}
}
