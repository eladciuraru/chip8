package main

import (
    "syscall"
    "unsafe"
)

type User32DLL struct {
    _DefWindowProcW   *syscall.Proc
    _RegisterClassW   *syscall.Proc
    _CreateWindowExW  *syscall.Proc
    _PeekMessageW     *syscall.Proc
    _TranslateMessage *syscall.Proc
    _DispatchMessageW *syscall.Proc
    _PostQuitMessage  *syscall.Proc
    _GetDC            *syscall.Proc
    _ReleaseDC        *syscall.Proc
    _AdjustWindowRect *syscall.Proc
    _GetClientRect    *syscall.Proc
    _GetKeyboardState *syscall.Proc
}


var user32 = NewUser32()


func NewUser32() *User32DLL {
    dll := syscall.MustLoadDLL("user32")

    return &User32DLL{
        _DefWindowProcW:   dll.MustFindProc("DefWindowProcW"),
        _RegisterClassW:   dll.MustFindProc("RegisterClassW"),
        _CreateWindowExW:  dll.MustFindProc("CreateWindowExW"),
        _PeekMessageW:     dll.MustFindProc("PeekMessageW"),
        _TranslateMessage: dll.MustFindProc("TranslateMessage"),
        _DispatchMessageW: dll.MustFindProc("DispatchMessageW"),
        _PostQuitMessage:  dll.MustFindProc("PostQuitMessage"),
        _GetDC:            dll.MustFindProc("GetDC"),
        _ReleaseDC:        dll.MustFindProc("ReleaseDC"),
        _AdjustWindowRect: dll.MustFindProc("AdjustWindowRect"),
        _GetClientRect:    dll.MustFindProc("GetClientRect"),
        _GetKeyboardState: dll.MustFindProc("GetKeyboardState"),
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


func (u32 *User32DLL) PeekMessageW(lpMsg *MSG, hWnd syscall.Handle,
                                  msgFilterMin, msgFilterMax, wRemoveMsg uint32) int32 {
    ret, _, _ := u32._PeekMessageW.Call(
        uintptr(unsafe.Pointer(lpMsg)),
        uintptr(hWnd),
        uintptr(msgFilterMin),
        uintptr(msgFilterMax),
        uintptr(wRemoveMsg),
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


func (u32 *User32DLL) PostQuitMessage(nExitCode int32) {
    u32._PostQuitMessage.Call(uintptr(nExitCode))
}


func (u32 *User32DLL) GetDC(hWnd syscall.Handle) syscall.Handle {
    ret, _, _ := u32._GetDC.Call(
        uintptr(hWnd),
    )

    return syscall.Handle(ret)
}


func (u32 *User32DLL) ReleaseDC(hWnd, hdc syscall.Handle) int32 {
    ret, _, _ := u32._ReleaseDC.Call(
        uintptr(hWnd),
        uintptr(hdc),
    )

    return int32(ret)
}


func (u32 *User32DLL) AdjustWindowRect(lpRect *RECT, dwStyle uint32,
                                       bMenu int32) int32 {
    ret, _, _ := u32._AdjustWindowRect.Call(
        uintptr(unsafe.Pointer(lpRect)),
        uintptr(dwStyle),
        uintptr(bMenu),
    )

    return int32(ret)
}


func (u32 *User32DLL) GetClientRect(hWnd syscall.Handle, lpRect *RECT) int32 {
    ret, _, _ := u32._GetClientRect.Call(
        uintptr(hWnd),
        uintptr(unsafe.Pointer(lpRect)),
    )

    return int32(ret)
}


func (u32 *User32DLL) GetKeyboardState(lpKeyState *byte) int32 {
    ret, _, _ := u32._GetKeyboardState.Call(
        uintptr(unsafe.Pointer(lpKeyState)),
    )

    return int32(ret)
}
