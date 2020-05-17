[![GoDoc](https://godoc.org/github.com/ryn1x/gtiff?status.svg)](https://godoc.org/github.com/ryn1x/gtiff)

# gtiff
gtiff provides simple reading and writing of uint8, uint16, and float32 grayscale tiff images.
Per the [TIFF 6.0 spec](https://www.adobe.io/content/dam/udp/en/open/standards/tiff/TIFF6.pdf) grayscale images are 4 or 8 bit, but 16 and 32 bit images are still common in scientific and medical imaging.
Although basic, this package provides functionality not found in other full-featured packages that strictly adhere to the spec.
This package currently only supports the minimum tags required per the spec and does not offer much to manipulate them.
New and more advanced features will be added as I personally need them. Pull requests are always welcome!

## Usage
The intended usage of this package is to read the data from a tiff image, manipulate the data as necessary, and write a new tiff image.

Example (error handling omitted for brevity):
```go
package main

import (
    "os"

    "github.com/ryn1x/gtiff"
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
```
## License
gtiff is available under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).
