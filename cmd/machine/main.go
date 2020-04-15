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

    // This is the mapping from CHIP8 input pad to platform keyboard
    //     1 2 3 C ----> 1 2 3 4
    //     4 5 6 D ----> Q W E R
    //     7 8 9 E ----> A S D F
    //     A 0 B F ----> Z X C V
    keyMap := map[byte]byte{
        // First row
        machine.KeyOne   : '1',
        machine.KeyTwo   : '2',
        machine.KeyThree : '3',
        machine.KeyC     : '4',

        // Second row
        machine.KeyFour : 'Q',
        machine.KeyFive : 'W',
        machine.KeySix  : 'E',
        machine.KeyD    : 'R',

        // Third row
        machine.KeySeven : 'A',
        machine.KeyEigth : 'S',
        machine.KeyNine  : 'D',
        machine.KeyE     : 'F',

        // Fourth row
        machine.KeyA    : 'Z',
        machine.KeyZero : 'X',
        machine.KeyB    : 'C',
        machine.KeyF    : 'V',
    }
    fmt.Println(keyMap[1])

    fillCell := func(bitmap *Bitmap, x, y, count uint, flag bool) {
        value := byte(0)
        if flag {
            value = byte(255)
        }

        var offX, offY uint
        for offY = 0; offY < count; offY++ {
            for offX = 0; offX < count; offX++ {
                offset := (offY + (y * count)) * uint(bitmap.stride) +
                          (offX + (x * count)) * uint(bitmap.pixelSize)

                bitmap.buffer[offset]     = value
                bitmap.buffer[offset + 1] = value
                bitmap.buffer[offset + 2] = value
                bitmap.buffer[offset + 3] = value
            }
        }
    }

    window.WindowLoop(func(win *Window) {
        vm.DoCycle()

        // From platform keyboard to CHIP8 input pad
        keyStates := vm.GetKeyStates()
        for key := range keyStates {
            keyStates[key] = win.keyboard[keyMap[byte(key)]]
        }
        // keyStates[machine.KeyOne] = true
        vm.SetKeyStates(keyStates)

        // win.bitmap.buffer = make([]byte, len(win.bitmap.buffer))

        // From CHIP8 graphics memory to window bitmap buffer
        for index, state := range vm.GetDisplayState() {
            baseY := uint(index) / uint(machine.DisplayWidth)
            baseX := uint(index) % uint(machine.DisplayWidth)
            fillCell(win.bitmap, baseX, baseY, args.scale, state)
        }
    })
}


func Assert(cond bool, err error) {
    if cond == false {
        panic(err)
    }
}
