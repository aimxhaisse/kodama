package main

import(
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"log"
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

// Filter creates a new image which is a filtered copy of the input
type Filter interface {
     Process(image.Image) image.Image
}

// Brightness is a filter that modifies the brightness of the image
type Brightness struct {
     Factor	uint32	// Percentage of brightness to apply
}

// Process applies a brightness filter to the image
func (filter *Brightness) Process(in image.Image) image.Image {
     bounds := in.Bounds()
     out := image.NewRGBA(bounds)
     for x := bounds.Min.X; x < bounds.Max.X; x++ {
          for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
	      r, g, b, a := in.At(x, y).RGBA()

	      nr := Trunc(r + (0xFFFF * filter.Factor) / 100)
	      ng := Trunc(g + (0xFFFF * filter.Factor) / 100)
	      nb := Trunc(b + (0xFFFF * filter.Factor) / 100)

	      nc := color.NRGBA64{nr, ng, nb, uint16(a)}
	      out.Set(x, y, nc)
	  }
     }
     return out.SubImage(bounds)
}

// NewBrightness creates a new filter for brightness
func NewBrightness(factor uint32) *Brightness {
     return &Brightness {
     	    factor,
     }
}

// Darkness is a filter that modifies the darkness of the image
type Darkness struct {
     Factor	uint32	// Percentage of darkness to apply
}

// Process applies a darkness filter to the image
func (filter *Darkness) Process(in image.Image) image.Image {
     bounds := in.Bounds()
     out := image.NewRGBA(bounds)
     for x := bounds.Min.X; x < bounds.Max.X; x++ {
          for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
	      r, g, b, a := in.At(x, y).RGBA()

	      nr := Strunc(int32(r - (0xFFFF * filter.Factor) / 100))
	      ng := Strunc(int32(g - (0xFFFF * filter.Factor) / 100))
	      nb := Strunc(int32(b - (0xFFFF * filter.Factor) / 100))

	      nc := color.NRGBA64{nr, ng, nb, uint16(a)}
	      out.Set(x, y, nc)
	  }
     }
     return out.SubImage(bounds)
}

// NewDarkness creates a new filter for darkness
func NewDarkness(factor uint32) *Darkness {
     return &Darkness {
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
     input := GetImage("sample.jpg")
     filter := NewDarkness(50)
     output := filter.Process(input)
     PutImage(output, "/home/mxs/vhost/www/shots/paris/processed.jpg")
}
