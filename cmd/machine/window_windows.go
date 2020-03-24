package main

import (
    "fmt"
    "syscall"
    "time"
)

type Window struct {
    width    int32
    height   int32
    clock    time.Duration
    bitmap   *Bitmap
    keyboard [256]bool

    // Windows shit
    hWnd syscall.Handle
}


// TODO: Chane this function to support optional pattern
func NewWindow(title string, width, height int32, clock time.Duration) *Window {
    window := &Window{
        width:  width,
        height: height,
        clock:  time.Second / clock,
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

    // Given width & height are the desired client area,
    // overall window size needs to be calculated
    style := WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU |
             WS_MINIMIZEBOX | WS_VISIBLE
    rect  := RECT{left: 0, top: 0, right: width, bottom: height}
    user32.AdjustWindowRect(&rect, style, 0)

    window.hWnd = user32.CreateWindowExW(0, lpTitle, lpTitle, style,
                                         int32(CW_USEDEFAULT), int32(CW_USEDEFAULT),
                                         rect.right - rect.left, rect.bottom - rect.top,
                                         0, 0, hInstance, 0)
    Assert(window.hWnd != 0, fmt.Errorf("failed to create window"))

    return window
}


func (win *Window) WindowProc(hWnd syscall.Handle, Msg uint32,
                              wParam uintptr, lParam uintptr) uintptr {
    switch Msg {
        // case WM_KEYDOWN: fallthrough
        // case WM_KEYUP:
        //     win.keyboard[wParam] = uint32(lParam) & (1 << 31) == 0

        case WM_DESTROY:
            user32.PostQuitMessage(0)
            return 0
    }

    return user32.DefWindowProcW(hWnd, Msg, wParam, lParam)
}


func (win *Window) FlushBitmap() {
    hdc := user32.GetDC(win.hWnd)

    bitmap := win.bitmap
    gdi32.StretchDIBits(hdc, 0, 0, win.width, win.height,
                        0, 0, bitmap.width, bitmap.height,
                        &bitmap.buffer[0], &bitmap.info, DIB_RGB_COLORS, SRCCOPY)

    user32.ReleaseDC(win.hWnd, hdc)
}


func (win *Window) MessageLoop() bool {
    msg := MSG{}
    for user32.PeekMessageW(&msg, 0, 0, 0, PM_REMOVE) != 0 {
        user32.TranslateMessage(&msg)
        user32.DispatchMessageW(&msg)

        if msg.message == WM_QUIT {
            return false
        }
    }

    return true
}


func (win *Window) HandleKeyboard() {
    var keyState [256]byte

    user32.GetKeyboardState(&keyState[0])

    for i, state := range keyState {
        win.keyboard[i] = (state & 0x80) != 0
    }
}


func (win *Window) WindowLoop(callback func(*Window)) {
    for win.MessageLoop() {
        start := time.Now()

        win.HandleKeyboard()
        callback(win)
        win.FlushBitmap()

        // Limit FPS
        time.Sleep(win.clock - time.Since(start))
    }
}
