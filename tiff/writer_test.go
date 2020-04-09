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
	err = WriteTiff8(w, data8, 5, 5)
	if err != nil {
		t.Errorf("test failed")
	}
}

var data8 = []uint8{
	255, 0, 0, 0, 255,
	0, 255, 0, 255, 0,
	0, 0, 255, 0, 0,
	0, 255, 0, 255, 0,
	255, 0, 0, 0, 255}
