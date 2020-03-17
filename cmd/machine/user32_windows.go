package main

import (
	"syscall"
	"unsafe"
)

type User32DLL struct {
    _DefWindowProcW   *syscall.Proc
    _RegisterClassW   *syscall.Proc
    _CreateWindowExW  *syscall.Proc
    _GetMessageW      *syscall.Proc
    _TranslateMessage *syscall.Proc
    _DispatchMessageW *syscall.Proc
}


var user32 = newUser32() 


func newUser32() *User32DLL {
    dll := syscall.MustLoadDLL("user32")

    return &User32DLL{
        _DefWindowProcW:   dll.MustFindProc("DefWindowProcW"),
        _RegisterClassW:   dll.MustFindProc("RegisterClassW"),
        _CreateWindowExW:  dll.MustFindProc("CreateWindowExW"),
        _GetMessageW:      dll.MustFindProc("GetMessageW"),
        _TranslateMessage: dll.MustFindProc("TranslateMessage"),
        _DispatchMessageW: dll.MustFindProc("DispatchMessageW"),
    }
}


func (u32 *User32DLL) DefWindowProcW(hWnd syscall.Handle, Msg uint32,
                                     wParam uintptr, lParam uintptr) uintptr {
    ret, _, _ := u32._DefWindowProcW.Call(
        uintptr(hWnd),
        uintptr(Msg),
        wParam,
        lParam,
    )

    return ret
}


func (u32 *User32DLL) RegisterClassW(lpWndClass *WNDCLASSW) uint16 {
    ret, _, _ := u32._RegisterClassW.Call(
        uintptr(unsafe.Pointer(lpWndClass)),
    )

    return uint16(ret)
}


func (u32 *User32DLL) CreateWindowExW(dwExStyle uint32, lpClassName, lpWindowName *uint16,
                                      dwStyle uint32, X, Y, nWidth, nHeight int32,
                                      hWndParent, hMenu, hInstance syscall.Handle,
                                      lpParam uintptr) syscall.Handle {
    ret, _, _ := u32._CreateWindowExW.Call(
        uintptr(dwExStyle),
        uintptr(unsafe.Pointer(lpClassName)),
        uintptr(unsafe.Pointer(lpWindowName)),
        uintptr(dwStyle),
        uintptr(X),
        uintptr(Y),
        uintptr(nWidth),
        uintptr(nHeight),
        uintptr(hMenu),
        uintptr(hWndParent),
        uintptr(hInstance),
        lpParam,
    )

    return syscall.Handle(ret)
}


func (u32 *User32DLL) GetMessageW(lpMsg *MSG, hWnd syscall.Handle, 
                                  msgFilterMin uint32, msgFilterMax uint32) int32 {
    ret, _, _ := u32._GetMessageW.Call(
        uintptr(unsafe.Pointer(lpMsg)),
        uintptr(hWnd),
        uintptr(msgFilterMin),
        uintptr(msgFilterMax),
    )

    return int32(ret)
}


func (u32 *User32DLL) TranslateMessage(lpMsg *MSG) bool {
    ret, _, _ := u32._TranslateMessage.Call(uintptr(unsafe.Pointer(lpMsg)))

    return ret != 0
}


func (u32 *User32DLL) DispatchMessageW(lpMsg *MSG) uintptr {
    ret, _, _ := u32._DispatchMessageW.Call(uintptr(unsafe.Pointer(lpMsg)))

    return ret
}
