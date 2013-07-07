package cr2

// This is a hacky implementation of a cr2 decoder. I own a Canon 650D
// and since there's no official specification of this format, it may not
// work with other cameras. If you have trouble, I'd be glad to receive
// samples from other cameras.
//
// Thanks a lot to: http://lclevy.free.fr/cr2

import (
	"errors"
	"image"
	"io"
)

// Header of a little endian cr2 file:
// $ hexdump -C -n 12 IMG_1135.CR2
// 49 49 2a 00 10 00 00 00 43 52 02 00
const cr2Header = "\x49\x49\x2a\x00\x10\x00\x00\x00\x43\x52\x02\x00"

// Decode reads a CR2 image from r and returns it as an image.Image.
// The type of Image returned depends on the PNG contents.
func Decode(r io.Reader) (image.Image, error) {
	return nil, errors.New("cr2: not yet implemented")
}

// DecodeConfig returns the color model and dimensions of a CR2 image without
// decoding the entire image
func DecodeConfig(r io.Reader) (image.Config, error) {
	return image.Config{}, errors.New("cr2: not yet implemented")
}

func init() {
	image.RegisterFormat("cr2", cr2Header, Decode, DecodeConfig)
}
