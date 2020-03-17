package main

import (
	"fmt"
	"syscall"
	"time"
)

type Window struct {
	width  int32
    height int32
    bitmap *Bitmap
}


func NewWindow(title string, width, height int32) *Window {
    window := &Window{
        width:  width,
        height: height,
        bitmap: NewBitmap(width, height),
    }

    hInstance    := kernel32.GetModuleHandleW(nil)
    lpTitle, err := syscall.UTF16PtrFromString(title)
    Assert(err == nil, err)

    wndclass := WNDCLASSW{
        lpfnWndProc:   syscall.NewCallback(window.WindowProc),
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

    return window
}


func (win *Window) WindowProc(hWnd syscall.Handle, Msg uint32,
                              wParam uintptr, lParam uintptr) uintptr {
    switch Msg {
        case WM_PAINT:
            paint  := &PAINTSTRUCT{}
            hdc    := user32.BeginPaint(hWnd, paint)
            
            bitmap := win.bitmap
            gdi32.StretchDIBits(hdc, 0, 0, win.width, win.height,
                                0, 0, bitmap.width, bitmap.height,
                                &bitmap.buffer[0], &bitmap.info, DIB_RGB_COLORS, SRCCOPY)

            user32.EndPaint(hWnd, paint)
            return 0

        case WM_DESTROY:
            user32.PostQuitMessage(0)
            return 0
    }

    return user32.DefWindowProcW(hWnd, Msg, wParam, lParam)
}


func (win *Window) MessageLoop() {
    msg := MSG{}
    for {
        if ret := user32.GetMessageW(&msg, 0, 0, 0); ret == 0 || ret == -1 {
            break
        }

        user32.TranslateMessage(&msg)
        user32.DispatchMessageW(&msg)

        // Limit 60 FPS
        time.Sleep(1000 / 60)
    }
}
