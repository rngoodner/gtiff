package gtiff

import (
	"encoding/binary"
	"io"
)

type header struct {
	byteOrder      uint16
	tiffIdentifier uint16
	ifdOffset      uint32
}

type dir16 struct {
	tag   uint16
	dtype uint16
	count uint32
	value uint16
	pad   uint16
}

type dir32 struct {
	tag   uint16
	dtype uint16
	count uint32
	value uint32
}

// WriteTiff8 writes a tiff from a slice of uint8 data.
func WriteTiff8(w io.WriteSeeker, byteOrder binary.ByteOrder, data []uint8, width uint32, length uint32) error {
	// steps:
	// 1) write all image data starting at offset 8, seek to next word boundry and save offset
	// 2) write 1 ifd of 11 directory entries
	// 3) write all 11 directory entries, 1 for each required tag
	// 4) point stripOffset to 8
	// 5) write 4 bytes of 0 to indicate last ifd
	// 6) write header at offset 0 with ifdOffset from saved after image data

	// 1)
	if _, err := w.Seek(8, 0); err != nil {
		return err
	}
	if err := binary.Write(w, byteOrder, data); err != nil {
		return err
	}
	afterData, _ := w.Seek(0, io.SeekCurrent)
	// seek to next work boundry
	afterData = afterData/8*8 + 8
	w.Seek(afterData, 0)

	// 2)
	if err := binary.Write(w, byteOrder, uint16(11)); err != nil {
		return err
	}
	// 3-4)
	// ImageWidth
	if err := binary.Write(w, byteOrder, newDir32(256, width)); err != nil {
		return err
	}
	// ImageLength
	if err := binary.Write(w, byteOrder, newDir32(257, length)); err != nil {
		return err
	}
	// BitsPerSample
	if err := binary.Write(w, byteOrder, newDir16(258, 8)); err != nil {
		return err
	}
	// Compression
	if err := binary.Write(w, byteOrder, newDir16(259, 1)); err != nil {
		return err
	}
	// PhotometricInterpretation
	if err := binary.Write(w, byteOrder, newDir16(262, 1)); err != nil {
		return err
	}
	// StripOffsets
	if err := binary.Write(w, byteOrder, newDir32(273, 8)); err != nil {
		return err
	}
	// RowsPerStrip
	if err := binary.Write(w, byteOrder, newDir32(278, length)); err != nil {
		return err
	}
	// StripByteCounts
	if err := binary.Write(w, byteOrder, newDir32(279, width*length)); err != nil {
		return err
	}
	// XResolution
	if err := binary.Write(w, byteOrder, newDir32(282, 0)); err != nil {
		return err
	}
	// YResolution
	if err := binary.Write(w, byteOrder, newDir32(283, 0)); err != nil {
		return err
	}
	// ResolutionUnit
	if err := binary.Write(w, byteOrder, newDir16(296, 0)); err != nil {
		return err
	}

	// 5)
	if err := binary.Write(w, byteOrder, []byte{0, 0, 0, 0}); err != nil {
		return err
	}

	// 6)
	var bo uint16 = 0X4949 // default to little endian
	if byteOrder == binary.BigEndian {
		bo = 0X4D4D // big endian code
	}
	// create header
	h := header{bo, 42, uint32(afterData)}
	if _, err := w.Seek(0, 0); err != nil {
		return err
	}
	if err := binary.Write(w, byteOrder, h); err != nil {
		return err
	}

	return nil
}

// WriteTiff16 writes a tiff from a slice of uint16 data.
func WriteTiff16(w io.WriteSeeker, byteOrder binary.ByteOrder, data []uint16, width uint32, length uint32) error {
	// steps:
	// 1) write all image data starting at offset 8, seek to next word boundry and save offset
	// 2) write 1 ifd of 11 directory entries
	// 3) write all 11 directory entries, 1 for each required tag
	// 4) point stripOffset to 8
	// 5) write 4 bytes of 0 to indicate last ifd
	// 6) write header at offset 0 with ifdOffset from saved after image data

	// 1)
	if _, err := w.Seek(8, 0); err != nil {
		return err
	}
	if err := binary.Write(w, byteOrder, data); err != nil {
		return err
	}
	afterData, _ := w.Seek(0, io.SeekCurrent)
	// seek to next work boundry
	afterData = afterData/8*8 + 8
	w.Seek(afterData, 0)

	// 2)
	if err := binary.Write(w, byteOrder, uint16(11)); err != nil {
		return err
	}
	// 3-4)
	// ImageWidth
	if err := binary.Write(w, byteOrder, newDir32(256, width)); err != nil {
		return err
	}
	// ImageLength
	if err := binary.Write(w, byteOrder, newDir32(257, length)); err != nil {
		return err
	}
	// BitsPerSample
	if err := binary.Write(w, byteOrder, newDir16(258, 16)); err != nil {
		return err
	}
	// Compression
	if err := binary.Write(w, byteOrder, newDir16(259, 1)); err != nil {
		return err
	}
	// PhotometricInterpretation
	if err := binary.Write(w, byteOrder, newDir16(262, 1)); err != nil {
		return err
	}
	// StripOffsets
	if err := binary.Write(w, byteOrder, newDir32(273, 8)); err != nil {
		return err
	}
	// RowsPerStrip
	if err := binary.Write(w, byteOrder, newDir32(278, length)); err != nil {
		return err
	}
	// StripByteCounts
	if err := binary.Write(w, byteOrder, newDir32(279, width*length)); err != nil {
		return err
	}
	// XResolution
	if err := binary.Write(w, byteOrder, newDir32(282, 0)); err != nil {
		return err
	}
	// YResolution
	if err := binary.Write(w, byteOrder, newDir32(283, 0)); err != nil {
		return err
	}
	// ResolutionUnit
	if err := binary.Write(w, byteOrder, newDir16(296, 0)); err != nil {
		return err
	}

	// 5)
	if err := binary.Write(w, byteOrder, []byte{0, 0, 0, 0}); err != nil {
		return err
	}

	// 6)
	var bo uint16 = 0X4949 // default to little endian
	if byteOrder == binary.BigEndian {
		bo = 0X4D4D // big endian code
	}
	// create header
	h := header{bo, 42, uint32(afterData)}
	if _, err := w.Seek(0, 0); err != nil {
		return err
	}
	if err := binary.Write(w, byteOrder, h); err != nil {
		return err
	}

	return nil
}

// WriteTiff32 write a tiff from a slice of float32 data.
func WriteTiff32(w io.WriteSeeker, byteOrder binary.ByteOrder, data []float32, width uint32, length uint32) error {
	// steps:
	// 1) write all image data starting at offset 8, seek to next word boundry and save offset
	// 2) write 1 ifd of 12 directory entries
	// 3) write all 12 directory entries, 1 for each required tag + sample format
	// 4) point stripOffset to 8
	// 5) write 4 bytes of 0 to indicate last ifd
	// 6) write header at offset 0 with ifdOffset from saved after image data

	// 1)
	if _, err := w.Seek(8, 0); err != nil {
		return err
	}
	if err := binary.Write(w, byteOrder, data); err != nil {
		return err
	}
	afterData, _ := w.Seek(0, io.SeekCurrent)
	// seek to next work boundry
	afterData = afterData/8*8 + 8
	w.Seek(afterData, 0)

	// 2)
	if err := binary.Write(w, byteOrder, uint16(12)); err != nil {
		return err
	}
	// 3-4)
	// ImageWidth
	if err := binary.Write(w, byteOrder, newDir32(256, width)); err != nil {
		return err
	}
	// ImageLength
	if err := binary.Write(w, byteOrder, newDir32(257, length)); err != nil {
		return err
	}
	// BitsPerSample
	if err := binary.Write(w, byteOrder, newDir16(258, 32)); err != nil {
		return err
	}
	// Compression
	if err := binary.Write(w, byteOrder, newDir16(259, 1)); err != nil {
		return err
	}
	// PhotometricInterpretation
	if err := binary.Write(w, byteOrder, newDir16(262, 1)); err != nil {
		return err
	}
	// StripOffsets
	if err := binary.Write(w, byteOrder, newDir32(273, 8)); err != nil {
		return err
	}
	// RowsPerStrip
	if err := binary.Write(w, byteOrder, newDir32(278, length)); err != nil {
		return err
	}
	// StripByteCounts
	if err := binary.Write(w, byteOrder, newDir32(279, width*length)); err != nil {
		return err
	}
	// XResolution
	if err := binary.Write(w, byteOrder, newDir32(282, 0)); err != nil {
		return err
	}
	// YResolution
	if err := binary.Write(w, byteOrder, newDir32(283, 0)); err != nil {
		return err
	}
	// ResolutionUnit
	if err := binary.Write(w, byteOrder, newDir16(296, 0)); err != nil {
		return err
	}
	// SampleFormat (3 is IEEE floating point, default without tag is 1 uint)
	if err := binary.Write(w, byteOrder, newDir16(339, 3)); err != nil {
		return err
	}

	// 5)
	if err := binary.Write(w, byteOrder, []byte{0, 0, 0, 0}); err != nil {
		return err
	}

	// 6)
	var bo uint16 = 0X4949 // default to little endian
	if byteOrder == binary.BigEndian {
		bo = 0X4D4D // big endian code
	}
	// create header
	h := header{bo, 42, uint32(afterData)}
	if _, err := w.Seek(0, 0); err != nil {
		return err
	}
	if err := binary.Write(w, byteOrder, h); err != nil {
		return err
	}

	return nil
}

func newDir16(tag uint16, value uint16) dir16 {
	return dir16{tag, 3, 1, value, 0}
}

func newDir32(tag uint16, value uint32) dir32 {
	return dir32{tag, 4, 1, value}
}
