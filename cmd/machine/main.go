// +build windows

package main

import (
    "os"
    "fmt"
    "runtime"
    "flag"
    "time"

    "github.com/eladciuraru/chip8/pkg/machine"
)

type Arguments struct {
    path  string
    base  uint
    scale uint
    clock uint
}


func parseArgs() Arguments {
    const (
		pathHelp  = "Path to ROM file contains CHIP-8 compiled instructions"
		baseHelp  = "The address to load the ROM file at"
		scaleHelp = "How much to scale each pixel"
		clockHelp = "What clock speed to use"
        shortHelp = " (shorthand)"
        
        scaleDefault = 10
        clockDefault = 60
    )

    var args Arguments

    flag.StringVar(&args.path, "path", args.path, pathHelp)
    flag.StringVar(&args.path, "p", args.path, pathHelp + shortHelp)

    flag.UintVar(&args.base, "base", args.base, baseHelp)
    flag.UintVar(&args.base, "b", args.base, baseHelp + shortHelp)

    flag.UintVar(&args.scale, "scale", scaleDefault, scaleHelp)
    flag.UintVar(&args.scale, "s", scaleDefault, scaleHelp + shortHelp)

    flag.UintVar(&args.clock, "clock", clockDefault, clockHelp)
    flag.UintVar(&args.clock, "c", clockDefault, clockHelp + shortHelp)

    flag.Parse()

    // Make path argument required
    if args.path == "" {
        flag.PrintDefaults()
        os.Exit(1)
    }

    return args
}


func main() {
    // This is required, since handling WM_PAINT message requires
    // the same thread that registered the WndProc
    runtime.LockOSThread()

    args := parseArgs()

    vm, err := machine.FromFile(args.path)
    if err != nil {
        panic(fmt.Errorf("failed to create chip8 virtual machine: %w", err))
    }

    // Setup the window
    width  := int32(machine.DisplayWidth) * int32(args.scale)
    height := int32(machine.DisplayHeight) * int32(args.scale)
    window := NewWindow("CHIP8 - Virtual Machine",
                        width, height, time.Duration(args.clock))

    index := uint(0)

    fillCell := func(bitmap *Bitmap, x, y, count uint) {
        var offX, offY uint
        for offY = 0; offY < count; offY++ {
            for offX = 0; offX < count; offX++ {
                offset := (offY + (y * count)) * uint(bitmap.stride) +
                          (offX + (x * count)) * uint(bitmap.pixelSize)

                bitmap.buffer[offset]     = 255
                bitmap.buffer[offset + 1] = 255
                bitmap.buffer[offset + 2] = 255
                bitmap.buffer[offset + 3] = 255
            }
        }
    }

    window.WindowLoop(func(win *Window) {
        vm.DoCycle()

        bitmap := win.bitmap
        bitmap.buffer = make([]byte, len(win.bitmap.buffer))
        
        baseY := uint(index) / uint(machine.DisplayWidth)
        baseX := uint(index) % uint(machine.DisplayWidth)
        fillCell(bitmap, baseX, baseY, args.scale)

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
        // index++
        index = index % 2048
    })
}


func Assert(cond bool, err error) {
    if cond == false {
        panic(err)
    }
}
