package main

import "unsafe"

type Bitmap struct {
    info      BITMAPINFO
    buffer    []byte
    width     int32
    height    int32
    stride    int32
    pixelSize int32
}


func NewBitmap(width, height int32) *Bitmap {
    const (
        bytesPerPixel = 4
        bitsPerByte   = 8
    )

    header := BITMAPINFOHEADER {
        biSize:        uint32(unsafe.Sizeof(BITMAPINFOHEADER{})),
        biWidth:       width,
        biHeight:      -height,  // For top down buffer
        biPlanes:      1,
        biBitCount:    bytesPerPixel * bitsPerByte,
        biCompression: BI_RGB,
    }

    stride := width * bytesPerPixel
    return &Bitmap{
        info:      BITMAPINFO{bmiHeader: header},
        width:     width,
        height:    height,
        stride:    stride,
        pixelSize: bytesPerPixel,
        buffer:    make([]byte, stride * height),
    }
}
