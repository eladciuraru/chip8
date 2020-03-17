package main

import (
	"fmt"
	"syscall"
)

type Window struct {
	width  int32
	height int32

	// Windows API needed shit
	hWnd syscall.Handle
}


func newWindow(title string, width, height int32) *Window {
    hInstance    := kernel32.GetModuleHandleW(nil)
    lpTitle, err := syscall.UTF16PtrFromString(title)
    Assert(err == nil, err)

    wndclass := WNDCLASSW{
        lpfnWndProc:   syscall.NewCallback(WindowProc),
        hInstance:     hInstance,
        lpszClassName: lpTitle,
    }
    
    res := user32.RegisterClassW(&wndclass)
    Assert(res != 0, fmt.Errorf("failed to register the window class"))

    style := WS_OVERLAPPED | WS_MINIMIZEBOX | WS_SYSMENU | WS_VISIBLE
    hWnd  := user32.CreateWindowExW(0, lpTitle, lpTitle, style,
                                    int32(CW_USEDEFAULT), int32(CW_USEDEFAULT),
                                    width, height, 0, 0, hInstance, 0)
    Assert(hWnd != 0, fmt.Errorf("failed to create window"))
    
    fmt.Println(wndclass)

    return &Window{
        width:  width,
        height: height,

        // Windows shit
        hWnd:   hWnd,
    }
}


func WindowProc(hWnd syscall.Handle, Msg uint32,
                wParam uintptr, lParam uintptr) uintptr {
    return user32.DefWindowProcW(hWnd, Msg, wParam, lParam)
}


func (win *Window) MessageLoop() {
    msg := MSG{}
    for user32.GetMessageW(&msg, 0, 0, 0) != 0 {
        // TODO: handle -1 returned from GetMessageW
        user32.TranslateMessage(&msg)
        user32.DispatchMessageW(&msg)
    }
}
