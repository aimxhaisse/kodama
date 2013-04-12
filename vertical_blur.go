package main

import (
	"image"
	"image/color"
)

// VerticalBlur is a filter that adds a vertical blur to the image
type VerticalBlur struct {
	Radius int
}

// Process applies a vertical blur filter to the image (efficient implementation)
func (filter *VerticalBlur) Process(in image.Image, out *image.RGBA, bounds image.Rectangle) {
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		prev_blur := filter.computeInitialBlur(in, bounds, x)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			prev_y := ClipInt(y-filter.Radius/2, 0, bounds.Max.Y-1)
			next_y := ClipInt(y+filter.Radius/2, 0, bounds.Max.Y-1)

			nb_elems := next_y - prev_y + 1

			pr, pg, pb, pa := in.At(x, prev_y).RGBA()
			nr, ng, nb, na := in.At(x, next_y).RGBA()
			vbr, vbg, vbb, vba := prev_blur.RGBA()

			cvbr := uint16(ClipInt((int(vbr)*nb_elems - int(pr) + int(nr)) / nb_elems, 0, 0xFFFF))
			cvbg := uint16(ClipInt((int(vbg)*nb_elems - int(pg) + int(ng)) / nb_elems, 0, 0xFFFF))
			cvbb := uint16(ClipInt((int(vbb)*nb_elems - int(pb) + int(nb)) / nb_elems, 0, 0xFFFF))
			cvba := uint16(ClipInt((int(vba)*nb_elems - int(pa) + int(na)) / nb_elems, 0, 0xFFFF))

			next_blur := color.NRGBA64{cvbr, cvbg, cvbb, cvba}
			out.Set(x, y, next_blur)
			prev_blur = next_blur
		}
	}
}

func (filter *VerticalBlur) computeInitialBlur(in image.Image, bounds image.Rectangle, x int) color.Color {
	start := ClipInt(bounds.Min.Y-filter.Radius/2, 0, bounds.Max.Y)
	end := ClipInt(bounds.Min.Y+filter.Radius/2, 0, bounds.Max.Y)

	var vbr, vbg, vbb, vba int
	for iter := start; iter <= end; iter++ {
		r, g, b, a := in.At(x, iter).RGBA()
		vbr += int(r)
		vbg += int(g)
		vbb += int(b)
		vba += int(a)
	}

	nb_iter := (end - start) + 1
	r := uint16(ClipInt(vbr / nb_iter, 0, 0xFFFF))
	g := uint16(ClipInt(vbg / nb_iter, 0, 0xFFFF))
	b := uint16(ClipInt(vbb / nb_iter, 0, 0xFFFF))
	a := uint16(ClipInt(vba / nb_iter, 0, 0xFFFF))

	return color.NRGBA64{r, g, b, a}
}

// NewVerticalBlur creates a new filter for blur
func NewVerticalBlur(radius int) *VerticalBlur {
	return &VerticalBlur{
		radius,
	}
}
