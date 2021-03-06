package filters

import (
	"errors"
	"fmt"
	"image/color"
	"strconv"
)

// Saturation is a filter that modifies the saturation of the image
type Saturation struct {
	Strength uint32 // Percentage of saturation to apply
}

// NewSaturation creates a new filter for darkness
func NewSaturation(argv []string) (*Saturation, error) {
	if len(argv) != 2 {
		return nil, errors.New("invalid syntax for saturation, expected usage: saturation <strength>")
	}
	strength, err := strconv.Atoi(argv[1])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid parameter for saturation: %s", err.Error()))
	}
	if strength > 0 {
		return &Saturation{
			uint32(strength),
		}, nil
	}
	return nil, errors.New("parameter 'strength' must be > 0")
}

// This filter is scalable
func (filter *Saturation) IsScalable() {
}

// Process applies a saturation filter to the image
func (filter *Saturation) Process(img *FilterImage) error {
	out := img.Image
	bounds := out.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := out.At(x, y).RGBA()

			grey := (r + g + b) / 3

			nr := Trunc(r + ((Abs(int32(r-grey)) * filter.Strength) / 100))
			ng := Trunc(g + ((Abs(int32(g-grey)) * filter.Strength) / 100))
			nb := Trunc(b + ((Abs(int32(b-grey)) * filter.Strength) / 100))

			nc := color.NRGBA64{nr, ng, nb, uint16(a)}
			out.Set(x, y, nc)
		}
	}
	return nil
}
