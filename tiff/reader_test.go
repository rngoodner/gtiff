package tiff

import (
	"os"
	"reflect"
	"testing"
)

func TestReadIntegration8(t *testing.T) {
	r, err := os.Open("../test-images/cell8.tif")
	if err != nil {
		t.Fatal(err)
	}

	// read tags
	tags, header, err := ReadTags(r)
	if err != nil {
		t.Fatal(err)
	}

	// read data
	data8, err := ReadData8(r, header, tags)
	if err != nil {
		t.Fatal(err)
	}

	// close tiff file
	r.Close()

	expected8 := []uint8{117, 119, 118, 117, 118, 119, 119, 119, 118, 122}
	if !reflect.DeepEqual(expected8, data8[30359:]) {
		t.Errorf("expected %v, got %v", expected8, data8[30359:])
	}
}

func TestReadIntegration16(t *testing.T) {
	r, err := os.Open("../test-images/cell16.tif")
	if err != nil {
		t.Fatal(err)
	}

	// read tags
	tags, header, err := ReadTags(r)
	if err != nil {
		t.Fatal(err)
	}

	// read data
	data16, err := ReadData16(r, header, tags)
	if err != nil {
		t.Fatal(err)
	}

	// close tiff file
	r.Close()

	expected16 := []uint16{34492, 35354, 34923, 34492, 34923, 35354, 35354, 35354, 34923, 36648}
	if !reflect.DeepEqual(expected16, data16[30359:]) {
		t.Errorf("expected %v, got %v", expected16, data16[30359:])
	}
}

func TestReadIntegration32(t *testing.T) {
	r, err := os.Open("../test-images/cell32.tif")
	if err != nil {
		t.Fatal(err)
	}

	// read tags
	tags, header, err := ReadTags(r)
	if err != nil {
		t.Fatal(err)
	}

	// read data
	data32, err := ReadData32(r, header, tags)
	if err != nil {
		t.Fatal(err)
	}

	// close tiff file
	r.Close()

	expected32 := []float32{0.48369575, 0.49456015, 0.48912793, 0.48369575, 0.48912793, 0.49456015, 0.49456015, 0.49456015, 0.48912793, 0.51087207}
	if !reflect.DeepEqual(expected32, data32[30359:]) {
		t.Errorf("expected %v, got %v", expected32, data32[30359:])
	}

}
