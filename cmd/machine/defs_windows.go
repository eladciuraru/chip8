package main

import "syscall"

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

// Using var instead of constants because go's constants
// are shit to deal with
var (
    WS_OVERLAPPED     uint32 = 0x00000000
	WS_MINIMIZEBOX    uint32 = 0x00020000
	WS_SYSMENU        uint32 = 0x00080000
    WS_VISIBLE        uint32 = 0x10000000

    CW_USEDEFAULT     uint32 = 0x80000000
)
