package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ryn1x/grayscale-tiff/tiff"
)

// print tags of tiff supplied via first command line args
func main() {
	r, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	// read tags
	tags, _, err := tiff.ReadTags(r)
	if err != nil {
		log.Fatal(err)
	}

	// print tags
	fmt.Println(tags)
}
