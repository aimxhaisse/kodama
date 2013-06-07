package filters

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"strconv"
)

// Darkness is a filter that modifies the darkness of the image
type Darkness struct {
	Strength uint32 // Percentage of darkness to apply
}

// NewDarkness creates a new filter for darkness
func NewDarkness(argv []string) (*Darkness, error) {
	if len(argv) != 2 {
		return nil, errors.New("invalid syntax for darkness, expected usage: darkness <strength>")
	}
	strength, err := strconv.Atoi(argv[1])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid parameter for darkness: %s", err.Error()))
	}
	if strength > 0 {
		return &Darkness{
			uint32(strength),
		}, nil
	}
	return nil, errors.New("parameter 'strength' must be > 0")
}

// This filter is scalable
func (filter *Darkness) IsScalable() {
}

// Process applies a darkness filter to the image
func (filter *Darkness) Process(img *FilterImage) error {
	in := img.Image
	bounds := in.Bounds()
	out := image.NewRGBA64(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := in.At(x, y).RGBA()

			nr := Strunc(int32(r - (0xFFFF*filter.Strength)/100))
			ng := Strunc(int32(g - (0xFFFF*filter.Strength)/100))
			nb := Strunc(int32(b - (0xFFFF*filter.Strength)/100))

			nc := color.NRGBA64{nr, ng, nb, uint16(a)}
			out.Set(x, y, nc)
		}
	}
	img.Image = out
	return nil
}
