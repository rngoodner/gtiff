package tiff

import (
	"log"
	"os"
	"testing"
)

func TestWrite8(t *testing.T) {
	fileName := "../test-images/test-output8.tif"
	w, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Could not open file: %v", fileName)
	}
	defer w.Close()
	WriteTiff8(w, data8, 5, 5)

	res := false
	if !res {
		t.Errorf("test failed")
	}
}

var data8 = []uint8{
	1, 0, 0, 0, 1,
	0, 1, 0, 1, 0,
	0, 0, 1, 0, 0,
	0, 1, 0, 1, 0,
	1, 0, 0, 0, 1}
