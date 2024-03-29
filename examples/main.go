package main

import (
	"os"

	"github.com/rngoodner/gtiff"
)

func main() {
	// open a tiff file
	r, _ := os.Open("../test-images/cell32.tif") // error handling omitted
	defer r.Close()

	// read tags
	tags, header, _ := gtiff.ReadTags(r) // error handling omitted

	// read data
	data, _ := gtiff.ReadData32(r, header, tags) // error handling omitted

	// >>> manipulate data as desired here <<<

	// write a new tiff
	w, _ := os.Create("../test-images/example-output-cell32.tif") // error handling omitted
	defer w.Close()
	gtiff.WriteTiff32(w, header.ByteOrder, data, tags.ImageWidth, tags.ImageLength) // error handling omitted
}
