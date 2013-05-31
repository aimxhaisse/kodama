package filters

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"strconv"
)

// Brightness is a filter that modifies the brightness of the image
type Brightness struct {
	Strength uint32 // Percentage of brightness to apply
}

// NewBrightness creates a new filter for brightness
func NewBrightness(argv []string) (*Brightness, error) {
	if len(argv) != 2 {
		return nil, errors.New("invalid syntax for brightness, expected usage: brightness <strength>")
	}
	strength, err := strconv.Atoi(argv[1])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid parameter for brightness: %s", err.Error()))
	}
	if strength > 0 {
		return &Brightness{
			uint32(strength),
		}, nil
	}
	return nil, errors.New("parameter 'strength' must be > 0")
}

// This filter is scalable
func (filter *Brightness) IsScalable() {
}

// Process applies a brightness filter to the image
func (filter *Brightness) Process(in image.Image, out *image.RGBA, bounds image.Rectangle) {
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := in.At(x, y).RGBA()

			// r, g, b, a are 16bits components in a uint32
			nr := Trunc(r + (0xFFFF*filter.Strength)/100)
			ng := Trunc(g + (0xFFFF*filter.Strength)/100)
			nb := Trunc(b + (0xFFFF*filter.Strength)/100)

			nc := color.NRGBA64{nr, ng, nb, uint16(a)}
			out.Set(x, y, nc)
		}
	}
}
