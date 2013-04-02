package main

import (
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(4)
	input := GetImage("sample.jpg")
	saturate := NewSaturation(115)
	output := ProcessInParallel(input, 4, saturate)
	PutImage(output, "processed.jpg")
}
