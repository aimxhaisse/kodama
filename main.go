package main

import (
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(4)
	input := GetImage("sample.jpg")
	blur := NewBlur(50)
	output := ProcessInParallel(input, 4, blur)
	PutImage(output, "processed.jpg")
}
