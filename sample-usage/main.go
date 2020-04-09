package main

import (
	"os"

	"github.com/ryn1x/grayscale-tiff/tiff"
)

func main() {
	// open a tiff file
	r, _ := os.Open("../test-images/cell32.tif") // error handling omitted
	defer r.Close()

	// read tags
	tags, header, _ := tiff.ReadTags(r) // error handling omitted

	// read data
	data, _ := tiff.ReadData32(r, header, tags) // error handling omitted

	// >>> manipulate data as desired here <<<

	// write out a new tiff
	fileName := "../test-images/sample-output-cell32.tif"
	w, _ := os.Create(fileName) // error handling omitted
	defer w.Close()
	tiff.WriteTiff32(w, data, tags.ImageWidth, tags.ImageLength) // error handling omitted
}
