package main

import "syscall"

// Using var instead of constants because go's constants
// are shit to deal with
var (
	WS_OVERLAPPED     uint32 = 0x00000000
	WS_CAPTION        uint32 = 0x00C00000
	WS_SYSMENU        uint32 = 0x00080000
	WS_MINIMIZEBOX    uint32 = 0x00020000
	WS_VISIBLE        uint32 = 0x10000000

	CW_USEDEFAULT     uint32 = 0x80000000
	
	WM_DESTROY        uint32 = 0x00000002
	WM_QUIT           uint32 = 0x00000012

	PM_REMOVE         uint32 = 0x00000001

	BI_RGB            uint32 = 0x00000000
	DIB_RGB_COLORS    uint32 = 0x00000000
	SRCCOPY           uint32 = 0x00CC0020
)


type WNDCLASSW struct {
	style         uint32
	lpfnWndProc   uintptr
	cbClsExtra    int32
	cbWndExtra    int32
	hInstance     syscall.Handle
	hIcon         syscall.Handle
	hCursor       syscall.Handle
	hbrBackground syscall.Handle
	lpszMenuName  *uint16
	lpszClassName *uint16
}


type POINT struct {
	x int32
	y int32
}


type MSG struct {
	hwnd     syscall.Handle
	message  uint32
	wParam	 uintptr
	lParam   uintptr
	time     uint32
	pt       POINT
	lPrivate uint32
}


type BITMAPINFOHEADER struct {
	biSize          uint32
	biWidth         int32
	biHeight        int32
	biPlanes        uint16
	biBitCount      uint16
	biCompression   uint32
	biSizeImage     uint32
	biXPelsPerMeter int32
	biYPelsPerMeter int32
	biClrUsed       uint32
	biClrImportant  uint32
}

type BITMAPINFO struct {
	bmiHeader BITMAPINFOHEADER
	bmiColors [1]uint32
}


type RECT struct {
	left   int32
	top    int32
	right  int32
	bottom int32
}


type PAINTSTRUCT struct {
	hdc         syscall.Handle
	fErase      int32
	rcPaint     RECT
	fRestore    int32
	fIncUpdate  int32
	rgbReserved [32]byte
}
