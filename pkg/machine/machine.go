package machine

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const (
    InstructionSize uint = 0x02
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
    memory   [MemorySize]byte
    video    [VideoWidth][VideoHeight]bool
    keyboard [16]bool
}


type Bus interface {
    Read(uint16) byte
    Write(uint16, byte)
}


func New(rom []byte) *VirtualMachine {
    vm := &VirtualMachine{}

    // Share the bus with the CPU
    vm.cpu = NewProcessor(vm)

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


func (vm *VirtualMachine) fixAddress(addr uint16) uint16 {
    // Address should only use 12 bits, so only look at lowest 12 bits.
    // This has the effect of looping around the address space in case of
    // address bigger than 12 bits
    return addr & 0x0FFF
}


func (vm *VirtualMachine) Read(addr uint16) byte {
    return vm.memory[vm.fixAddress(addr)]
}


func (vm *VirtualMachine) Write(addr uint16, data byte) {
    vm.memory[vm.fixAddress(addr)] = data
}


func (vm *VirtualMachine) Start() {
    for i := 0; i < 0x11; i++ {
        vm.cpu.Cycle()
    }

    fmt.Println("Finished")
}
