package gtiff

import (
	"encoding/binary"
	"fmt"
)

// Header represents a  parsed tiff header.
type Header struct {
	ByteOrder      binary.ByteOrder
	TiffIdentifier uint16
	IFDOffset      uint32
}

// Tags holds the minumum grayscale tag set per tiff 6.0 spec.
type Tags struct {
	ImageWidth                uint32   // 256 (short or long)
	ImageLength               uint32   // 257 (short or long)
	BitsPerSample             uint16   // 258 (count: single value for grayscale images)
	Compression               uint16   // 259
	PhotometricInterpretation uint16   // 262
	StripOffsets              []uint32 // 273 (short or long) (count: StripsPerImage)
	RowsPerStrip              uint32   // 278 (short or long)
	StripByteCounts           []uint32 // 279 (short or long) (count: StripsPerImage)
	XResolution               []uint32 // 282 (count: 2, numerator, denomenator)
	YResolution               []uint32 // 283 (count: 2, numerator, denomenator)
	ResolutionUnit            uint16   // 296
}

// String method for Tags
func (t Tags) String() string {
	res := ""
	res += fmt.Sprintf("ImageWidth(256):                %v\n", t.ImageWidth)
	res += fmt.Sprintf("Imagelength(257):               %v\n", t.ImageLength)
	res += fmt.Sprintf("BitsPerSample(258):             %v\n", t.BitsPerSample)
	res += fmt.Sprintf("Compression(259):               %v\n", t.Compression)
	res += fmt.Sprintf("PhotometricInterpretation(262): %v\n", t.PhotometricInterpretation)
	res += fmt.Sprintf("StripOffsets(273):              %v\n", t.StripOffsets)
	res += fmt.Sprintf("RowsPerStrip(278):              %v\n", t.RowsPerStrip)
	res += fmt.Sprintf("StripByteCounts(279):           %v\n", t.StripByteCounts)
	res += fmt.Sprintf("XResolution(282):               %v\n", t.XResolution)
	res += fmt.Sprintf("YResolution(283):               %v\n", t.YResolution)
	res += fmt.Sprintf("ResolutionUnit(296):            %v", t.ResolutionUnit)
	return res
}
