package tiff

import (
	"io"
)

func WriteTiff8(w io.WriteSeeker, data []uint8, width uint32, length uint32) {
	// create header
	h := struct {
		byteOrder      uint16
		tiffIdentifier uint16
		ifdOffset      uint32
	}{0X4D4D, 42, 64}

	// create tags directory entries

	// ImageWidth
	// ImageLength
	// BitsPerSample
	// Compression
	// PhotometricInterpretation
	// StripOffsets
	// RowsPerStrip
	// StripByteCounts
	// XResolution
	// YResolution
	// ResolutionUnit

	// write image

	w.Write(data)
}
