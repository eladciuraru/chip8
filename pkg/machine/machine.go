package machine

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const (
    InstructionSize uint = 0x02

    // Memory map related consts
    MemorySize         uint16 = 0x1000
    MemoryRomAddr      uint16 = 0x0200
    MemoryStackAddr    uint16 = 0x0EA0
    MemoryWorkAreaAddr uint16 = 0x0ED0
    MemoryKeyboardAddr uint16 = 0x0EF0
    MemoryDisplayAddr  uint16 = 0x0F00

    // Display realted consts
    DisplayWidth      uint = 64
    DisplayHeight     uint = 32
    DisplaySize       uint = DisplayHeight * DisplayWidth
    DisplayMemorySize uint = uint(MemorySize - MemoryDisplayAddr)
    DisplayPixelWidth uint = 8
)

type VirtualMachine struct {
    cpu      *Processor
    memory   [MemorySize]byte
    display  [DisplaySize]bool
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
    romMaxSize := MemoryStackAddr - MemoryRomAddr
    limitCopy  := minUint16(uint16(len(rom)), romMaxSize)
    copy(vm.memory[MemoryRomAddr:], rom[:limitCopy])

    // TODO: Load sprites data

    return vm
}


func FromReader(reader io.Reader) (*VirtualMachine, error) {
    romData, err := ioutil.ReadAll(reader)
    if err != nil {
        return nil, fmt.Errorf("failed to read all: %w", err)
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
    var value byte
    switch addr := vm.fixAddress(addr); {
        // Handle keyboard
        case MemoryKeyboardAddr <= addr && addr < MemoryDisplayAddr:
            index := addr - MemoryKeyboardAddr
            if vm.keyboard[index] {
                value = 1  // Key is pressesd
            }

        // Regular RAM access - keyboard address range doesn't 
        // need special handling for reading, 
        // but I rather it being expressed explicitly 
        case MemoryDisplayAddr <= addr && addr < MemorySize: fallthrough
        default: value = vm.memory[addr]
    }

    return value
}


func (vm *VirtualMachine) Write(addr uint16, data byte) {
    addr = vm.fixAddress(addr)
    switch {
        // Handle keyboard - Copy changes to the keyboard bool buffer
        case MemoryKeyboardAddr <= addr && addr < MemoryDisplayAddr:
            index := addr - MemoryKeyboardAddr
            vm.keyboard[index] = data == 1

        // Regular RAM access - keyboard address range doesn't 
        // need special handling for reading, 
        // but I rather it being expressed explicitly 
        
        // Handle display - Copy changes to the display bool buffer
        // This will be easier for the user of this pacakge to handle,
        // since this is a monochrome display
        case MemoryDisplayAddr <= addr && addr < MemorySize:
            index := uint(addr - MemoryDisplayAddr) / DisplayPixelWidth
            data  := data
            for offset := uint(0); offset < DisplayPixelWidth; offset++ {
                vm.display[index + offset] = data & 1 == 1

                data >>= 1
            }
    }

    vm.memory[addr] = data
}


func (vm *VirtualMachine) Start() {

}
