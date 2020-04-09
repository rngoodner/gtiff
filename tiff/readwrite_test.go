package tiff

import (
	"os"
	"testing"
)

func TestReadWrite32(t *testing.T) {
	r, err := os.Open("../test-images/cell32.tif")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

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

	// write back out data
	fileName := "../test-images/test-output-cell32.tif"
	w, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("Could not open file: %v", fileName)
	}
	defer w.Close()
	err = WriteTiff32(w, data32, tags.ImageWidth, tags.ImageLength)
	if err != nil {
		t.Fatal(err)
	}
}
