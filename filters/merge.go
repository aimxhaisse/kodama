package filters

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"os"
)

// Merge is a filter that merge the current image with the input image
type Merge struct {
	Image image.Image  // input image
}

// NewMerge creates a new merge filter
func NewMerge(argv[] string) (*Merge, error) {
	if len(argv) != 2 {
		return nil, errors.New("invalid syntax for merge, expected usage: merge <input>")
	}

	reader, err := os.Open(argv[1])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't open input file: %s", err.Error()))
	}
	defer reader.Close()
	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't decode input file: %s", err.Error()))
	}

	return &Merge{
		m,
	}, nil
}

// Process merges the input image
func (filter *Merge) Process(img *FilterImage) error {
	out := img.Image
	bounds := out.Bounds()
	inbounds := filter.Image.Bounds()

	xmin := bounds.Min.X
	if xmin < inbounds.Min.X {
		xmin = inbounds.Min.X
	}
	xmax := bounds.Max.X
	if xmax > inbounds.Max.X {
		xmax = inbounds.Max.X
	}
	ymin := bounds.Min.Y
	if ymin < inbounds.Min.Y {
		ymin = inbounds.Min.Y
	}
	ymax := bounds.Max.Y
	if ymax > inbounds.Max.Y {
		ymax = inbounds.Max.Y
	}

	for x := xmin; x < xmax; x++ {
		for y := ymin; y < ymax; y++ {
			r, g, b, a := out.At(x, y).RGBA()
			ir, ig, ib, ia := filter.Image.At(x, y).RGBA()

			r = uint32(ClipInt(int(r + ir), 0, 0xFFFF))
			g = uint32(ClipInt(int(g + ig), 0, 0xFFFF))
			b = uint32(ClipInt(int(b + ib), 0, 0xFFFF))
			a = uint32(ClipInt(int(a + ia), 0, 0xFFFF))

			out.Set(x, y, color.NRGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
		}
	}

	return nil
}
