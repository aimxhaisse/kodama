package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"runtime"
)

// Trunc truncates a color component to a 16bits value
func Trunc(component uint32) uint16 {
	if component > 0xFFFF {
		return 0xFFFF
	}
	return uint16(component)
}

// Strunc truncates a signed color component to a 16bits value
func Strunc(component int32) uint16 {
	if component > 0xFFFF {
		return 0xFFFF
	}
	if component < 0 {
		return 0
	}
	return uint16(component)
}

// Abs returns the absolute value of in
func Abs(in int32) uint32 {
	if in > 0 {
		return uint32(in)
	}
	return uint32(-in)
}

// Filter creates a new image which is a filtered copy of the input
type Filter interface {
	Process(in image.Image, out *image.RGBA, area image.Rectangle)
}

// Brightness is a filter that modifies the brightness of the image
type Brightness struct {
	Factor uint32 // Percentage of brightness to apply
}

// Process applies a brightness filter to the image
func (filter *Brightness) Process(in image.Image, out *image.RGBA, bounds image.Rectangle) {
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := in.At(x, y).RGBA()

			nr := Trunc(r + (0xFFFF*filter.Factor)/100)
			ng := Trunc(g + (0xFFFF*filter.Factor)/100)
			nb := Trunc(b + (0xFFFF*filter.Factor)/100)

			nc := color.NRGBA64{nr, ng, nb, uint16(a)}
			out.Set(x, y, nc)
		}
	}
}

// NewBrightness creates a new filter for brightness
func NewBrightness(factor uint32) *Brightness {
	return &Brightness{
		factor,
	}
}

// Darkness is a filter that modifies the darkness of the image
type Darkness struct {
	Factor uint32 // Percentage of darkness to apply
}

// Process applies a darkness filter to the image
func (filter *Darkness) Process(in image.Image, out *image.RGBA, bounds image.Rectangle) {
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := in.At(x, y).RGBA()

			nr := Strunc(int32(r - (0xFFFF * filter.Factor)/100))
			ng := Strunc(int32(g - (0xFFFF * filter.Factor)/100))
			nb := Strunc(int32(b - (0xFFFF * filter.Factor)/100))

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

// Saturation is a filter that modifies the saturation of the image
type Saturation struct {
	Factor uint32 // Percentage of saturation to apply
}

// Process applies a darkness filter to the image
func (filter *Saturation) Process(in image.Image, out *image.RGBA, bounds image.Rectangle) {
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := in.At(x, y).RGBA()

			grey := (r + g + b) / 3

			nr := Trunc(r + ((Abs(int32(r - grey)) * filter.Factor) / 100))
			ng := Trunc(g + ((Abs(int32(g - grey)) * filter.Factor) / 100))
			nb := Trunc(b + ((Abs(int32(b - grey)) * filter.Factor) / 100))

			nc := color.NRGBA64{nr, ng, nb, uint16(a)}
			out.Set(x, y, nc)
		}
	}
}

// Split a job and send chunks to several workers
func ProcessInParallel(in image.Image, jobs int, worker Filter) image.Image {
	bounds := in.Bounds()
	x_unit := (bounds.Max.X - bounds.Min.X) / jobs
	ch := make(chan bool)
	out := image.NewRGBA(bounds)
	for i := 0; i < jobs; i++ {
		min := image.Point{i * x_unit, bounds.Min.Y}
		max := image.Point{(i + 1) * x_unit, bounds.Max.Y}
		area := image.Rectangle{min, max}
		log.Printf("first job going from [%d,%d] to [%d,%d]", area.Min.X, area.Min.Y, area.Max.X, area.Max.Y)
		go func (ch chan bool) {
			worker.Process(in, out, area)
			ch <- true
		}(ch)
	}
	done := 0
	// Wait for workers to complete
	for done != jobs {
		select {
		case <- ch:
			done++
		}
	}
	return out.SubImage(bounds)
}

// NewSaturation creates a new filter for darkness
func NewSaturation(factor uint32) *Saturation {
	return &Saturation{
		factor,
	}
}

// GetImage returns the image pointed by path
func GetImage(path string) image.Image {
	file, err := os.Open("sample.jpg")
	if err != nil {
		log.Fatal(err)
	}
	image, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	return image
}

// PutImage write the image to path
func PutImage(image image.Image, path string) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = jpeg.Encode(file, image, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	runtime.GOMAXPROCS(4)
	input := GetImage("sample.jpg")

	brigthness := NewBrightness(40)
	darkness := NewDarkness(20)
	saturate := NewSaturation(115)

	output := ProcessInParallel(ProcessInParallel(ProcessInParallel(input, 4, darkness), 4, saturate), 4, brigthness)

	PutImage(output, "processed.jpg")
}
