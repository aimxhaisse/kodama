package filters

import "image"

// Filters is a wrapper around images
type FilterImage struct {
	Image *image.RGBA64
}

// Filter processes an image
type Filter interface {
	Process(img *FilterImage) error
}

// ClipInt clips an integer between min and max
func ClipInt(i int, min int, max int) int {
	if i > max {
		return max
	}
	if i < min {
		return min
	}
	return i
}

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
