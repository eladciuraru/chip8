package machine

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)


const (
    MemorySize      uint = 0x1000
    RomMaxSize      uint = 0x0D00
    RomStartAddress uint = 0x0200
    VideoWidth      uint = 64
    VideoHeight     uint = 32
)


// Memory Map:
// 0x0000 - 0x01FF : Reserved for interpreter
// 0x0200 - 0x0E9F : Program / Data space
// 0x0EA0 - 0x0EFF : Stack space
// 0x0F00 - 0x0FFF : Display buffer
type VirtualMachine struct {
    cpu      *Processor
    bus      Bus
    memory   [MemorySize]byte
    video    [VideoWidth][VideoHeight]bool
    keyboard [16]bool
}


func New(rom []byte) *VirtualMachine {
    vm := &VirtualMachine{}

    // Share the bus with the CPU
    vm.cpu = NewProcessor(vm.bus, uint16(RomStartAddress))

    // Load rom data into ram
    limitCopy := minUint(uint(len(rom)), RomMaxSize)
    copy(vm.memory[RomStartAddress:], rom[:limitCopy])
    // Later copy font data into ram

    return vm
}


func FromReader(reader io.Reader) (*VirtualMachine, error) {
    romData, err := ioutil.ReadAll(reader)
    if err != nil {
        return nil, fmt.Errorf("failed to read all: %w", err)
    } else if uint(len(romData)) > RomMaxSize {
        return nil, fmt.Errorf("rom size: %d exceed maximum allowed size: %d",
                               len(romData), RomMaxSize)
    }

    return New(romData), nil
}


func FromFile(filename string) (*VirtualMachine, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open %v: %w", filename, err)
    }
    defer file.Close()

    return FromReader(file)
}


// Why isn't there a builtin for this, seems stupid that a 
// high level language like go has many things like this missing
func minUint(a, b uint) uint {
    if a <= b {
        return a
    } else {
        return b
    }
}
