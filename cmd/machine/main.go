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
    
    index  := 0
    window := NewWindow("test", 640, 320)
    window.WindowLoop(func(win *Window) {
        start := time.Now()

        if index >= 2048 {
            index = 0
        }
        win.bitmap.buffer = make([]byte, len(win.bitmap.buffer))

        stride := int(machine.DisplayWidth * 4 * 10)
        var x, y int
        for y = 0; y < 10; y++ {
            for x = 0; x < 10; x++ {
                offset := (y + (index / 64)*10)* stride + (x + (index % 64) * 10) * 4
                win.bitmap.buffer[offset]     = 255
                win.bitmap.buffer[offset + 1] = 255
                win.bitmap.buffer[offset + 2] = 255
            }
        }
        fmt.Println(index)
        index++

        // Limit 60 FPS
        duration := (time.Second / 60) - time.Since(start)
        time.Sleep(duration)
    })

    // vm.SetKeyState(machine.KeyA, true)
}


func Assert(cond bool, err error) {
	if cond == false {
		panic(err)
	}
}
