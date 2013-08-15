package cr2

// This is a hacky implementation of a cr2 decoder. I own a Canon 650D
// and since there's no official specification of this format, it may not
// work with other cameras. If you have trouble, I'd be glad to receive
// samples.
//
// Thanks a lot to: http://lclevy.free.fr/cr2

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"io"
	"io/ioutil"
	"fmt"
	"bufio"
)

// Header of a little endian cr2 file:
// $ hexdump -C -n 12 IMG_1135.CR2
// 49 49 2a 00 10 00 00 00 43 52 02 00
const cr2Header = "\x49\x49\x2a\x00????\x43\x52\x02\x00"

// decoder is the internal representation of a cr2 file
type decoder struct {
	buf *bytes.Reader      // entire cr2 file in memory
}

// scanHeader scans the tiff header
func (d decoder) scanHeader() error {
	head := make([]byte, 4)
	n, err := d.buf.Read(head)
	if err != nil {
		return nil
	}
	if n != 4 || string(head) != "\x49\x49\x2a\x00" {
		return errors.New("cr2: not a cr2 file")
	}
	return nil
}

// tiffTag represents a tiff tag
type tiffTag struct {
	Id    uint16
	Kind  uint16
	Nb    uint32
	Value uint32
}

// available kinds for a tiffTag
const (
	KIND_UNKNOWN = iota
	KIND_UCHAR
	KIND_STRING
	KIND_USHORT
	KIND_ULONG
	KIND_URATIO
	KIND_CHAR
	KIND_BYTES
	KIND_SHORT
	KIND_LONG
	KIND_RATIO
	KIND_FLOAT32
	KIND_FLOAT64
)

// known tiff tags names
type TiffTagNames map[uint16]string

var tiff_tag_names = TiffTagNames{
	0x0100: "imageWidth",
	0x0101: "imageHeight",
	0x0102: "bitsPerSample",
	0x0103: "compression",
	0x010f: "maker",
	0x0110: "model",
	0x0111: "stripOffset",
	0x0112: "orientaton",
	0x0117: "stripByteCounts",
	0x011a: "xResolution",
	0x011b: "yResolution",
	0x0128: "resolutionUnit",
	0x0132: "dateTime",
	0x8769: "EXIF",
	0x8825: "GPSData",
}

// String dumps the attributes of a tiffTag
func (t tiffTag) String() string {
	return fmt.Sprintf("id=0x%04x, kind=0x%04x, nb=0x%08x", t.Id, t.Kind, t.Nb)
}

// prettify returns a key and a value representing the tag
func (t tiffTag) prettify(r *bytes.Reader) (k string, v string, err error) {
	k, ok := tiff_tag_names[t.Id]
	if !ok {
		k = fmt.Sprintf("unknownTag(0x%04X)", t.Id)
	}

	// backup current offset before reading in buffer
	back, err := r.Seek(0, 1)
	if err != nil {
		return k, v, err
	}

	// ensure offset is restored
	defer func (back int64, r *bytes.Reader) {
		r.Seek(back, 0)
	}(back, r)

	switch (t.Kind) {

	case KIND_STRING:
		// move to string location and scan string
		_, err = r.Seek(int64(t.Value), 0)
		if err == nil {
			s := bufio.NewReader(r)
			v, err = s.ReadString(0)
			if err == nil && len(v) > 0 {
				v = v[:len(v) - 1]
			}
		}

		return k, v, err

	default:
		v = fmt.Sprintf("%d", t.Value);
		return k, v, err
		
	}

	return k, "unknownValue", err
}

// scanImageFileEntry scans a tiff tag
func (d decoder) scanImageFileEntry() error {
	var tag tiffTag
	err := binary.Read(d.buf, binary.LittleEndian, &tag)
	if err != nil {
		return err
	}
	k, v, err := tag.prettify(d.buf)
	if err != nil {
		return err
	}
	fmt.Printf("%s: %s\n", k, v)
	return nil
}

// scanImageFileDirectory scans all tiff tags in an IFD
func (d decoder) scanImageFileDirectory() error {
	// get the number of entries of the IFD and skip those
	var nb_entries uint16
	err := binary.Read(d.buf, binary.LittleEndian, &nb_entries)
	if err != nil {
		return err
	}
	for i := uint16(0); i < nb_entries; i++ {
		err = d.scanImageFileEntry()
		if err != nil {
			return err
		}
	}

	return nil
}

// Decode reads a CR2 image from r and returns it as an image.Image.
// The type of Image returned depends on the PNG contents.
func Decode(r io.Reader) (image.Image, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	d := &decoder{bytes.NewReader(buf)}
	err = d.scanHeader()
	if err != nil {
		return nil, err
	}

	// get the offset of the first IFD section and move there
	var offset uint32
	err = binary.Read(d.buf, binary.LittleEndian, &offset)
	if err != nil {
		return nil, err
	}
	_, err = d.buf.Seek(int64(offset), 0)
	if err != nil {
		return nil, err
	}

	// CR2 format includes 4 sections, each is composed of a
	// header containing metadata and a picture.
	//
	// We only deal with the fourth picture, which has the highest
	// resolution (others are thumbnails designed for camera use).
	for i := 0; i < 4; i++ {
		err = d.scanImageFileDirectory()
		if err != nil {
			return nil, err
		}
		err = binary.Read(d.buf, binary.LittleEndian, &offset)
		if err != nil {
			return nil, err
		}
		_, err = d.buf.Seek(int64(offset), 0)
		if err != nil {
			return nil, err
		}
	}

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