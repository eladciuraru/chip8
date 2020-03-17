package main

import "unsafe"

type Bitmap struct {
    info   BITMAPINFO
    buffer []byte
    width  int32
    height int32
}


func NewBitmap(width, height int32) *Bitmap {
    const (
        bytesPerPixel = 4
        bitsPerByte   = 8
    )

    header := BITMAPINFOHEADER {
        biSize:        uint32(unsafe.Sizeof(BITMAPINFOHEADER{})),
        biWidth:       width,
        biHeight:      height,
        biPlanes:      1,
        biBitCount:    bytesPerPixel * bitsPerByte,
        biCompression: BI_RGB,
    }

    return &Bitmap{
        info:   BITMAPINFO{bmiHeader: header},
        width:  width,
        height: height,
        buffer: make([]byte, width * height * bytesPerPixel),
    }
}
