package tiff

import (
	"os"
	"testing"
)

func TestWrite8(t *testing.T) {
	fileName := "../test-images/test-output8.tif"
	w, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("Could not open file: %v", fileName)
	}
	defer w.Close()
	err = WriteTiff8(w, data8, 5, 5)
	if err != nil {
		t.Errorf("test failed")
	}
}

func TestWrite16(t *testing.T) {
	fileName := "../test-images/test-output16.tif"
	w, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("Could not open file: %v", fileName)
	}
	defer w.Close()
	err = WriteTiff16(w, data16, 5, 5)
	if err != nil {
		t.Errorf("test failed")
	}
}

func TestWrite32(t *testing.T) {
	fileName := "../test-images/test-output32.tif"
	w, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("Could not open file: %v", fileName)
	}
	defer w.Close()
	err = WriteTiff32(w, data32, 5, 5)
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

var data16 = []uint16{
	65535, 0, 0, 0, 65535,
	0, 65535, 0, 65535, 0,
	0, 0, 65535, 0, 0,
	0, 65535, 0, 65535, 0,
	65535, 0, 0, 0, 65535}

var data32 = []float32{
	1.0, 0.0, 0.0, 0.0, 1.0,
	0.0, 1.0, 0.0, 1.0, 0.0,
	0.0, 0.0, 1.0, 0.0, 0.0,
	0.0, 1.0, 0.0, 1.0, 0.0,
	1.0, 0.0, 0.0, 0.0, 1.0}
