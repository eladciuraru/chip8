// +build windows

package main

import (
    "fmt"
    "runtime"
    "time"

    "github.com/eladciuraru/chip8/pkg/machine"
)

func main() {
    // This is required, since handling WM_PAINT message requires
    // the same thread that registered the WndProc
    runtime.LockOSThread()

    index  := uint(0)
    window := NewWindow("test", 640, 320)
    window.WindowLoop(func(win *Window) {
        start := time.Now()

        win.bitmap.buffer = make([]byte, len(win.bitmap.buffer))

        stride := uint(machine.DisplayWidth * 4 * 10)
        var x, y uint
        for y = 0; y < 10; y++ {
            for x = 0; x < 10; x++ {
                offset := (y + (index / 64)*10)* stride + (x + (index % 64) * 10) * 4
                win.bitmap.buffer[offset]     = 255
                win.bitmap.buffer[offset + 1] = 255
                win.bitmap.buffer[offset + 2] = 255
            }
        }
        fmt.Println(index)
        if win.keyboard['W'] {
            index -= 64  // up
        }
        if win.keyboard['A'] {
            index -= 1  // left
        }
        if win.keyboard['S'] {
            index += 64  // down
        }
        if win.keyboard['D'] {
            index += 1  // right
        }
        index = index % 2048
        // index++

        // Limit 60 FPS
        duration := (time.Second / 60) - time.Since(start)
        time.Sleep(duration)
    })
}


func Assert(cond bool, err error) {
    if cond == false {
        panic(err)
    }
}
