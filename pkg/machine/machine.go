package machine

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const (
    InstructionSize uint16 = 2
    FontSpriteSize  uint16 = 5
)


const (
    // Memory map related consts
    MemorySize          uint16 = 0x1000
    MemoryFontTableAddr uint16 = 0x0000
    MemoryRomAddr       uint16 = 0x0200
    MemoryStackAddr     uint16 = 0x0EA0
    MemoryWorkAreaAddr  uint16 = 0x0ED0
    MemoryKeyboardAddr  uint16 = 0x0EF0
    MemoryDisplayAddr   uint16 = 0x0F00
)


const (
    // Display realted consts
    DisplayWidth        uint16 = 64
    DisplayHeight       uint16 = 32
    DisplaySize         uint16 = DisplayHeight * DisplayWidth
    DisplayMemorySize   uint16 = MemorySize - MemoryDisplayAddr
    DisplayPixelWidth   uint16 = 8
    DisplayMemoryStride uint16 = DisplayWidth / DisplayPixelWidth
)


const (
    // Keyboard related consts
    KeyZero          byte = iota
    KeyOne
    KeyTwo
    KeyThree
    KeyFour
    KeyFive
    KeySix
    KeySeven
    KeyEigth
    KeyNine
    KeyA
    KeyB
    KeyC
    KeyD
    KeyE
    KeyF

    KeyboardSize     byte = 16
    KeyPressedValue  byte = 1
    KeyReleasedValue byte = 0
)

type VirtualMachine struct {
    cpu      *Processor
    memory   [MemorySize]byte
    display  [DisplaySize]bool
    keyboard [KeyboardSize]bool
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

    // Load sprites data - currently only fonts
    copy(vm.memory[MemoryFontTableAddr:], fontTable)

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
                value = KeyPressedValue
            } else {
                value = KeyReleasedValue
            }

        // Handle display
        case MemoryDisplayAddr <= addr && addr < MemorySize:
            offset := (addr - MemoryDisplayAddr) * DisplayPixelWidth
            mask   := byte(0x80)
            for i := uint16(0); i < DisplayPixelWidth; i++ {
                if vm.display[offset + i] {
                    value |= mask
                }
                mask >>= 1
            }

        // Regular RAM access
        default:
            value = vm.memory[addr]
    }

    return value
}


func (vm *VirtualMachine) Write(addr uint16, data byte) {
    switch addr = vm.fixAddress(addr); {
        // Handle keyboard - Copy changes to the keyboard bool buffer
        case MemoryKeyboardAddr <= addr && addr < MemoryDisplayAddr:
            index := addr - MemoryKeyboardAddr
            vm.keyboard[index] = (data == KeyPressedValue)

        // Handle display - Copy changes to the display bool buffer
        case MemoryDisplayAddr <= addr && addr < MemorySize:
            offset := (addr - MemoryDisplayAddr) * DisplayPixelWidth
            mask   := byte(0x80)
            for i := uint16(0); i < DisplayPixelWidth; i++ {
                vm.display[offset + i] = data & mask == mask

                mask >>= 1
            }

        // Regular RAM access
        default:
            vm.memory[addr] = data
    }
}


func (vm *VirtualMachine) DoCycle() {
    vm.cpu.Cycle()
}


func (vm *VirtualMachine) SetKeyState(keyIndex byte, pressed bool) bool {
    if keyIndex >= KeyboardSize {
        return false
    }

    prevState := vm.keyboard[keyIndex]
    vm.keyboard[keyIndex] = pressed

    return prevState
}


func (vm *VirtualMachine) GetKeyState(keyIndex byte) bool {
    if keyIndex >= KeyboardSize {
        return vm.keyboard[keyIndex]
    }

    return false
}


func (vm *VirtualMachine) GetKeyStates() [KeyboardSize]bool {
    return vm.keyboard
}


func (vm *VirtualMachine) SetKeyStates(keyStates [KeyboardSize]bool) [KeyboardSize]bool {
    var prevStates [KeyboardSize]bool

    copy(prevStates[:], vm.keyboard[:])
    copy(vm.keyboard[:], keyStates[:])

    return prevStates
}


func (vm *VirtualMachine) GetDisplayState() [DisplaySize]bool {
    return vm.display
}
