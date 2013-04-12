package main

import (
	"image"
	"image/jpeg"
	"log"
	"os"
)

// Filter creates a new image which is a filtered copy of the input
type Filter interface {
	Process(in image.Image, out *image.RGBA, area image.Rectangle)
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
		go func(ch chan bool) {
			worker.Process(in, out, area)
			ch <- true
		}(ch)
	}
	done := 0
	// Wait for workers to complete
	for done != jobs {
		select {
		case <-ch:
			done++
		}
	}
	return out.SubImage(bounds)
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
