package main

import (
	"syscall"
	"unsafe"
)

type GDI32DLL struct {
    _StretchDIBits      *syscall.Proc
}


var gdi32 = NewGDI32()


func NewGDI32() *GDI32DLL {
    dll := syscall.MustLoadDLL("gdi32")

    return &GDI32DLL{
        _StretchDIBits:      dll.MustFindProc("StretchDIBits"),
    }
}


func (g32 *GDI32DLL) StretchDIBits(hdc syscall.Handle, xDest, yDest, DestWidth, DestHeight,
                                   xSrc, ySrc, SrcWidth, SrcHeight int32, lpBits *byte, 
                                   lpbmi *BITMAPINFO, iUsage uint32, rop uint32) int32 {
    ret, _, _ := g32._StretchDIBits.Call(
        uintptr(hdc),
		uintptr(xDest),
		uintptr(yDest),
		uintptr(DestWidth),
		uintptr(DestHeight),
		uintptr(xSrc),
		uintptr(ySrc),
		uintptr(SrcWidth),
		uintptr(SrcHeight),
		uintptr(unsafe.Pointer(lpBits)),
		uintptr(unsafe.Pointer(lpbmi)),
		uintptr(iUsage),
		uintptr(rop),
    )

    return int32(ret)
}
