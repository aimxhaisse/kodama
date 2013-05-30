package filters

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"strconv"
)

// VBlur is a filter that adds a vertical blur to the image
type VBlur struct {
	Strength int
}

// NewVBlur creates a new filter for blur
func NewVBlur(argv []string) (*VBlur, error) {
	if len(argv) != 2 {
		return nil, errors.New("invalid syntax for vblur, expected usage: vblur <strength>")
	}
	strength, err := strconv.Atoi(argv[1])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid parameter for vblur: %s", err.Error()))
	}
	if strength > 0 {
		return &VBlur{
			strength,
		}, nil
	}
	return nil, errors.New("parameter 'strenght' must be > 0")
}

// IsScalable returns false because VBlur is not a scalable filter
func (filter *VBlur) IsScalable() bool {
	return false
}

// Process applies a vertical blur filter to the image (efficient implementation)
func (filter *VBlur) Process(in image.Image, out *image.RGBA, bounds image.Rectangle) {
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		prev_blur := filter.computeInitialBlur(in, bounds, x)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			prev_y := ClipInt(y-filter.Strength/2, 0, bounds.Max.Y-1)
			next_y := ClipInt(y+filter.Strength/2, 0, bounds.Max.Y-1)

			nb_elems := next_y - prev_y + 1

			pr, pg, pb, pa := in.At(x, prev_y).RGBA()
			nr, ng, nb, na := in.At(x, next_y).RGBA()
			vbr, vbg, vbb, vba := prev_blur.RGBA()

			cvbr := uint16(ClipInt((int(vbr)*nb_elems-int(pr)+int(nr))/nb_elems, 0, 0xFFFF))
			cvbg := uint16(ClipInt((int(vbg)*nb_elems-int(pg)+int(ng))/nb_elems, 0, 0xFFFF))
			cvbb := uint16(ClipInt((int(vbb)*nb_elems-int(pb)+int(nb))/nb_elems, 0, 0xFFFF))
			cvba := uint16(ClipInt((int(vba)*nb_elems-int(pa)+int(na))/nb_elems, 0, 0xFFFF))

			next_blur := color.NRGBA64{cvbr, cvbg, cvbb, cvba}
			out.Set(x, y, next_blur)
			prev_blur = next_blur
		}
	}
}

// computeInitialBlur computes the blur of the bound pixel
func (filter *VBlur) computeInitialBlur(in image.Image, bounds image.Rectangle, x int) color.Color {
	start := ClipInt(bounds.Min.Y-filter.Strength/2, 0, bounds.Max.Y)
	end := ClipInt(bounds.Min.Y+filter.Strength/2, 0, bounds.Max.Y)

	var vbr, vbg, vbb, vba int
	for iter := start; iter <= end; iter++ {
		r, g, b, a := in.At(x, iter).RGBA()
		vbr += int(r)
		vbg += int(g)
		vbb += int(b)
		vba += int(a)
	}

	nb_iter := (end - start) + 1
	r := uint16(ClipInt(vbr/nb_iter, 0, 0xFFFF))
	g := uint16(ClipInt(vbg/nb_iter, 0, 0xFFFF))
	b := uint16(ClipInt(vbb/nb_iter, 0, 0xFFFF))
	a := uint16(ClipInt(vba/nb_iter, 0, 0xFFFF))

	return color.NRGBA64{r, g, b, a}
}
