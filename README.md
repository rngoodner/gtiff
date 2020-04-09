# grayscale-tiff
grayscale-tiff provides simple reading and writing of uint8, uint16, and float32 grayscale tiff images.
Per the TIFF 6.0 spec (https://www.adobe.io/content/dam/udp/en/open/standards/tiff/TIFF6.pdf) grayscale images are 4 or 8 bit, but 16 and 32 bit images are still common in scientific and medical imaging.
Although basic, this package provides functionality not found in other full-featured packages that strictly adhere to the spec.
This package currently only supports the minimum tags required per the spec and does not offer much to manipulate them.
Data is currently hard-coded to write out in little-endian.
New and more advanced features will be added as I personally need them, but pull requests are welcome!

## Usage
The intended usage of this package is to read the data from a tiff image, manipulate the data as necessary, and write a new tiff image.

Example (error handling omitted for brevity):
```
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
```
## License
grayscale-tiff is available under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).
