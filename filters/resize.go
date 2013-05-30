package filters

import (
	"errors"
	"fmt"
	"image"
	"strconv"
)

// Resize is a filter that resizes the input image
type Resize struct {
	Width  int // new width
	Height int // new height
}

// NewResize creates a new filter for brightness
func NewResize(argv []string) (*Resize, error) {
	if len(argv) != 3 {
		return nil, errors.New("invalid syntax for resize, expected usage: resize <width> <height>")
	}
	w, err := strconv.Atoi(argv[1])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid parameter for width: %s", err.Error()))
	}
	h, err := strconv.Atoi(argv[1])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid parameter for height: %s", err.Error()))
	}

	if w <= 0 {
		return nil, errors.New("parameter 'width' must be > 0")
	}
	if h <= 0 {
		return nil, errors.New("parameter 'height' must be > 0")
	}

	return &Resize{
		w,
		h,
	}, nil
}

// IsScalable returns false because Resize is not a scalable filter
func (filter *Resize) IsScalable() bool {
	return false
}

// Process resizes the input image
func (filter *Resize) Process(in image.Image, out *image.RGBA, bounds image.Rectangle) {
	*out = *image.NewRGBA(image.Rect(0, 0, filter.Width, filter.Height))

	ratio_x := float64(bounds.Max.X) / float64(filter.Width)
	ratio_y := float64(bounds.Max.Y) / float64(filter.Height)

	for x := 0; x < filter.Width; x++ {
		for y := 0; y < filter.Height; y++ {
			in_x := int(ratio_x * float64(x))
			in_y := int(ratio_y * float64(y))
			out.Set(x, y, in.At(in_x, in_y))
		}
	}
}
