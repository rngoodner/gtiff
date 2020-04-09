// gtiff provides simple reading and writing of uint8, uint16, and float32 grayscale tiff images
package gtiff

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

// structure of a Directory Entry
type directoryEntry struct {
	Tag         uint16 // tag id number
	DType       uint16 // type of value
	Count       uint32 // number of values
	ValueOffset uint32 // offset to value
}

// reads the header of a Tiff file
func ReadHeader(r io.Reader) (Header, error) {
	var header Header
	header.ByteOrder = binary.BigEndian

	var byteOrder uint16

	// read byte order
	err := binary.Read(r, binary.BigEndian, &byteOrder)
	if err != nil {
		return header, err
	}

	// parse byte order
	switch byteOrder {
	case 0X4949:
		header.ByteOrder = binary.LittleEndian
	case 0X4D4D:
		header.ByteOrder = binary.BigEndian
	default:
		return header, errors.New("parse: invalid byte order")
	}

	// read tiff identifier order
	err = binary.Read(r, header.ByteOrder, &header.TiffIdentifier)
	if err != nil {
		return header, err
	}
	if header.TiffIdentifier != 42 {
		return header, fmt.Errorf("parse: invalid tiff identifier, expected: 42, got: %d", header.TiffIdentifier)
	}

	// read offset to first IFD
	err = binary.Read(r, header.ByteOrder, &header.IFDOffset)
	if err != nil {
		return header, err
	}

	return header, nil
}

// read all tags in the tiff file and record the values of supported tags
func ReadTags(r io.ReadSeeker) (Tags, Header, error) {
	var tags Tags

	header, err := ReadHeader(r)
	if err != nil {
		return tags, header, err
	}

	// offset to next IFD
	nextIFD := header.IFDOffset
	if _, err = r.Seek(int64(nextIFD), 0); err != nil {
		return tags, header, err
	}

	for nextIFD != 0 {
		// number of directory entries
		var numDE uint16
		err = binary.Read(r, header.ByteOrder, &numDE)
		if err != nil {
			return tags, header, err
		}

		// for each data directory
		var nextDir int64
		for i := uint16(0); i < numDE; i++ {
			// read static parts of directory entry
			var de directoryEntry
			err = binary.Read(r, header.ByteOrder, &de)
			if err != nil {
				return tags, header, err
			}

			// data type * number of values in bytes
			typeBytes16, _ := typeToBytes(de.DType)
			typeBytes := uint32(typeBytes16)
			typeBytes *= de.Count // bytes * number of values

			// if <= 4 bytes read value, else follow pointer to value
			if typeBytes <= 4 {
				// set directory entry value offset to current location in file
				offset, _ := r.Seek(0, io.SeekCurrent) // get current position in file
				de.ValueOffset = uint32(offset) - 4    // where we are now minus size of value offset (32bits=4bytes)
			}

			nextDir, _ = r.Seek(0, io.SeekCurrent) // get current position in file

			// if tag is supported then get the value(s), otherwise skip
			switch de.Tag {
			case 256:
				err = getTagValue16or32(r, &tags.ImageWidth, header.ByteOrder, de)
			case 257:
				err = getTagValue16or32(r, &tags.ImageLength, header.ByteOrder, de)
			case 258:
				err = getTagValue16(r, &tags.BitsPerSample, header.ByteOrder, de)
			case 259:
				err = getTagValue16(r, &tags.Compression, header.ByteOrder, de)
			case 262:
				err = getTagValue16(r, &tags.PhotometricInterpretation, header.ByteOrder, de)
			case 273:
				err = getMultiTagValues16or32(r, &tags.StripOffsets, header.ByteOrder, de)
			case 278:
				err = getTagValue16or32(r, &tags.RowsPerStrip, header.ByteOrder, de)
			case 279:
				err = getMultiTagValues16or32(r, &tags.StripByteCounts, header.ByteOrder, de)
			case 282:
				err = getMultiTagValues16or32(r, &tags.XResolution, header.ByteOrder, de)
			case 283:
				err = getMultiTagValues16or32(r, &tags.YResolution, header.ByteOrder, de)
			case 296:
				err = getTagValue16(r, &tags.ResolutionUnit, header.ByteOrder, de)
			default:
				continue
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: unable to read value for tag %d -- %s\n", de.Tag, err)
			}

			// seek to next dir
			if _, err = r.Seek(nextDir, 0); err != nil {
				return tags, header, err
			}
		}

		// get offset to next ifd
		err = binary.Read(r, header.ByteOrder, &nextIFD)
		if err != nil {
			return tags, header, err
		}
	}

	return tags, header, nil
}

// read 8 bit tiff image into a 1d slice
func ReadData8(r io.ReadSeeker, h Header, t Tags) ([]uint8, error) {
	var data []uint8
	for i, offset := range t.StripOffsets {
		// seek r to offset
		if _, err := r.Seek(int64(offset), 0); err != nil {
			return data, err
		}

		// read into slice
		numVals := t.StripByteCounts[i] / (uint32(t.BitsPerSample) / 8)
		stripData := make([]uint8, numVals)
		if err := binary.Read(r, h.ByteOrder, &stripData); err != nil {
			return data, err
		}
		data = append(data, stripData...)
	}
	return data, nil
}

// read 16 bit tiff image into a 1d slice
func ReadData16(r io.ReadSeeker, h Header, t Tags) ([]uint16, error) {
	var data []uint16
	for i, offset := range t.StripOffsets {
		// seek r to offset
		if _, err := r.Seek(int64(offset), 0); err != nil {
			return data, err
		}

		// read into slice
		numVals := t.StripByteCounts[i] / (uint32(t.BitsPerSample) / 8)
		stripData := make([]uint16, numVals)
		if err := binary.Read(r, h.ByteOrder, &stripData); err != nil {
			return data, err
		}
		data = append(data, stripData...)
	}
	return data, nil
}

// read 32 bit float tiff image into a 1d slice
func ReadData32(r io.ReadSeeker, h Header, t Tags) ([]float32, error) {
	var data []float32
	for i, offset := range t.StripOffsets {
		// seek r to offset
		if _, err := r.Seek(int64(offset), 0); err != nil {
			return data, err
		}

		// read into slice
		numVals := t.StripByteCounts[i] / (uint32(t.BitsPerSample) / 8)
		stripData := make([]float32, numVals)
		if err := binary.Read(r, h.ByteOrder, &stripData); err != nil {
			return data, err
		}
		data = append(data, stripData...)
	}
	return data, nil
}

// get value of a uint16 tag
func getTagValue16(r io.ReadSeeker, p *uint16, byteOrder binary.ByteOrder, de directoryEntry) error {
	if _, err := r.Seek(int64(de.ValueOffset), 0); err != nil {
		return err
	}

	if err := binary.Read(r, byteOrder, p); err != nil {
		return err
	}

	return nil
}

// get value of a uint32 tag
func getTagValue32(r io.ReadSeeker, p *uint32, byteOrder binary.ByteOrder, de directoryEntry) error {
	if _, err := r.Seek(int64(de.ValueOffset), 0); err != nil {
		return err
	}

	if err := binary.Read(r, byteOrder, p); err != nil {
		return err
	}

	return nil
}

// reads uint16 or uint32 value depending on type specified in directory entay and always return a uint32
func getTagValue16or32(r io.ReadSeeker, p *uint32, byteOrder binary.ByteOrder, de directoryEntry) error {
	var val16 uint16

	if _, err := r.Seek(int64(de.ValueOffset), 0); err != nil {
		return err
	}

	var err error
	switch de.DType {
	case 3:
		err = getTagValue16(r, &val16, byteOrder, de)
		*p = uint32(val16)
	case 4:
		err = getTagValue32(r, p, byteOrder, de)
	}
	if err != nil {
		return err
	}

	return nil
}

// populate slice with multiple values, reads uint16 or uint32 depending on type specified in directory entry and always returns uint32
func getMultiTagValues16or32(r io.ReadSeeker, p *[]uint32, byteOrder binary.ByteOrder, de directoryEntry) error {
	var newVal uint32

	for i := uint32(0); i < de.Count; i++ {
		if err := getTagValue16or32(r, &newVal, byteOrder, de); err != nil {
			return err
		}
		*p = append(*p, newVal)

		next, _ := typeToBytes(de.DType)
		de.ValueOffset += uint32(next)
	}

	return nil
}

// convert tiff numeric type to bytes
func typeToBytes(t uint16) (uint16, error) {
	// based on data type
	// if <= 4 bytes read value, else follow pointer to value
	var typeBytes uint16
	var err error
	switch t {
	case 1:
		typeBytes = 1 // byte
	case 2:
		typeBytes = 1 // ascii
	case 3:
		typeBytes = 2 // short
	case 4:
		typeBytes = 4 // long
	case 5:
		typeBytes = 8 // rational
	default:
		err = fmt.Errorf("type not supported, got %d, expected [1,5]", t)
	}
	return typeBytes, err
}
