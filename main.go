package main

import (
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(4)
	input := GetImage("sample.jpg")
	blur := NewHorizontalBlur(20)
	output := ProcessInParallel(input, 4, blur)
	PutImage(output, "processed.jpg")
}
