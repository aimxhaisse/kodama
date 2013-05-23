package filters

import (
	"image"
	"image/color"
)

// Blur is a filter that adds a blur to the image
type Blur struct {
	Radius int
}

// IsScalable returns true as this filter is scalable
func (filter *Blur) IsScalable() bool {
	return true
}

// Process applies a blur filter to the image
func (filter *Blur) Process(in image.Image, out *image.RGBA, bounds image.Rectangle) {
	// This is a naive implementation with a high complexity.
	// Each output pixel is the average of all pixels in its
	// surrounding box, thus complexity is W*H*R^2
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			_, _, _, a := in.At(x, y).RGBA()

			x_start := 0
			y_start := 0
			x_end := bounds.Max.X
			y_end := bounds.Max.Y

			if x-filter.Radius > 0 {
				x_start = x - filter.Radius
			}
			if x+filter.Radius < bounds.Max.X {
				x_end = x + filter.Radius
			}
			if y-filter.Radius > 0 {
				y_start = y - filter.Radius
			}
			if y+filter.Radius < bounds.Max.Y {
				y_end = y + filter.Radius
			}

			avg_r := uint32(0)
			avg_g := uint32(0)
			avg_b := uint32(0)
			pixels := uint32(0)
			for in_x := x_start; in_x <= x_end; in_x++ {
				for in_y := y_start; in_y <= y_end; in_y++ {
					in_r, in_g, in_b, _ := in.At(in_x, in_y).RGBA()
					avg_r += in_r
					avg_g += in_g
					avg_b += in_b
					pixels++
				}
			}

			nc := color.NRGBA64{uint16(avg_r / pixels), uint16(avg_g / pixels), uint16(avg_b / pixels), uint16(a)}
			out.Set(x, y, nc)
		}
	}
}

// NewBlur creates a new filter for blur
func NewBlur(radius int) *Blur {
	return &Blur{
		radius,
	}
}
