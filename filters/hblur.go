package filters

import (
	"image"
	"image/color"
)

// HBlur is a filter that adds a horizontal blur to the image
type HBlur struct {
	Radius int
}

// NewHBlur creates a new filter for horizontal blur
func NewHBlur(radius int) *HBlur {
	return &HBlur{
		radius,
	}
}

// IsScalable returns false because this filter is not scalable
func (filter *HBlur) IsScalable() bool {
	return false
}

// Process applies a horizontal blur filter to the image (efficient implementation)
func (filter *HBlur) Process(in image.Image, out *image.RGBA, bounds image.Rectangle) {
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		prev_blur := filter.computeInitialBlur(in, bounds, y)
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			prev_x := ClipInt(x-filter.Radius/2, 0, bounds.Max.X-1)
			next_x := ClipInt(x+filter.Radius/2, 0, bounds.Max.X-1)

			nb_elems := next_x - prev_x + 1

			pr, pg, pb, pa := in.At(prev_x, y).RGBA()
			nr, ng, nb, na := in.At(next_x, y).RGBA()
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
func (filter *HBlur) computeInitialBlur(in image.Image, bounds image.Rectangle, y int) color.Color {
	start := ClipInt(bounds.Min.X-filter.Radius/2, 0, bounds.Max.X)
	end := ClipInt(bounds.Min.X+filter.Radius/2, 0, bounds.Max.X)

	var vbr, vbg, vbb, vba int
	for iter := start; iter <= end; iter++ {
		r, g, b, a := in.At(iter, y).RGBA()
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
