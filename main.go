package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

// parsed tiff header
type Header struct {
	ByteOrder      binary.ByteOrder
	TiffIdentifier uint16
	IFDOffset      uint32
}

// minumum grayscale tag set per tiff 6.0 spec
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

// structure of a Directory Entry
type DirectoryEntry struct {
	Tag         uint16 // tag id number
	Type        uint16 // type of value
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
func ReadTags(r io.ReadSeeker) (Tags, error) {
	var tags Tags

	header, err := ReadHeader(r)
	if err != nil {
		return tags, err
	}

	// TODO add loop for avery IFD
	// ends when offset is 4 bytes of 0

	// offset to next IFD
	nextIFD := header.IFDOffset
	if _, err = r.Seek(int64(nextIFD), 0); err != nil {
		return tags, err
	}

	for nextIFD != 0 {
		// number of directory entries
		var numDE uint16
		err = binary.Read(r, header.ByteOrder, &numDE)
		if err != nil {
			return tags, err
		}

		//fmt.Printf("header: %v\n", header)
		//fmt.Printf("offset: %v\n", header.IFDOffset)
		//fmt.Printf("num de: %v\n", numDE)

		// for each data directory
		var nextDir int64
		for i := 0; i < int(numDE); i++ {
			// read static parts of directory entry
			var de DirectoryEntry
			err = binary.Read(r, header.ByteOrder, &de)
			if err != nil {
				return tags, err
			}

			//fmt.Printf("de%d: %v\n", i+1, de)

			// data type * number of values in bytes
			typeBytes16, _ := typeToBytes(de.Type)
			typeBytes := uint32(typeBytes16)
			typeBytes *= de.Count // bytes * number of values

			// if <= 4 bytes read value, else follow pointer to value
			if typeBytes <= 4 {
				// set directory entry value offset to current location in file
				offset, _ := r.Seek(0, io.SeekCurrent) // get current position in file
				de.ValueOffset = uint32(offset) - 4    // where we are now minus size of value offset (32bits=4bytes)
				//fmt.Printf("modified de%d: %v\n", i+1, de)
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
				return tags, err
			}
		}

		// get offset to next ifd
		err = binary.Read(r, header.ByteOrder, &nextIFD)
		if err != nil {
			return tags, err
		}
	}

	return tags, nil
}

// get value of a uint16 tag
func getTagValue16(r io.ReadSeeker, p *uint16, byteOrder binary.ByteOrder, de DirectoryEntry) error {
	if _, err := r.Seek(int64(de.ValueOffset), 0); err != nil {
		return err
	}

	if err := binary.Read(r, byteOrder, p); err != nil {
		return err
	}

	return nil
}

// get value of a uint32 tag
func getTagValue32(r io.ReadSeeker, p *uint32, byteOrder binary.ByteOrder, de DirectoryEntry) error {
	if _, err := r.Seek(int64(de.ValueOffset), 0); err != nil {
		return err
	}

	if err := binary.Read(r, byteOrder, p); err != nil {
		return err
	}

	return nil
}

// reads uint16 or uint32 value depending on type specified in directory entay and always return a uint32
func getTagValue16or32(r io.ReadSeeker, p *uint32, byteOrder binary.ByteOrder, de DirectoryEntry) error {
	var val16 uint16

	if _, err := r.Seek(int64(de.ValueOffset), 0); err != nil {
		return err
	}

	var err error
	switch de.Type {
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
func getMultiTagValues16or32(r io.ReadSeeker, p *[]uint32, byteOrder binary.ByteOrder, de DirectoryEntry) error {
	var newVal uint32

	for i := 0; i < int(de.Count); i++ {
		if err := getTagValue16or32(r, &newVal, byteOrder, de); err != nil {
			return err
		}
		*p = append(*p, newVal)

		next, _ := typeToBytes(de.Type)
		de.ValueOffset += uint32(next)
	}

	return nil
}

// contvert tiff numeric type to bytes
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

func main() {
	// open tiff file
	r, err := os.Open("cell32.tif")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	// read tags
	tags, err := ReadTags(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tags)

	// read data
}
