package disasm

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/eladciuraru/chip8/pkg/machine"
)

const InstructionSize uint = machine.InstructionSize


type Disasm struct {
    romData  []byte
    baseAddr uint
    labels   bool
    format   string
}


type Option func(*Disasm)


func WithAddress(addr uint) Option {
    return func(dis *Disasm) {
        dis.baseAddr = addr
    }
}


func WithLabels(dis *Disasm) {
    dis.labels = true
}


func WithFormat(format string) Option {
    return func(dis *Disasm) {
        dis.format = format
    }
}


func New(rom []byte, options ...Option) *Disasm {
    const (
		defaultFormat  = "addr mne op1 op2"
		defaultAddress = 0x200
    )

    disasm := &Disasm{
        baseAddr: defaultAddress,
        format:   defaultFormat,
        romData:  rom,
    }

    // This pattern was taken from https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
    for _, option := range options {
        option(disasm)
    }

    return disasm
}


func FromReader(reader io.Reader, options ...Option) (*Disasm, error) {
    romData, err := ioutil.ReadAll(reader)
    if err != nil {
        return nil, fmt.Errorf("failed to read all: %w", err)
    }

    return New(romData, options...), nil
}


func FromFile(filename string, options ...Option) (*Disasm, error) {
    // Maybe it is better to just use ioutil.ReadFile instead of reusing
    // `FromReader` function
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open %v: %w", filename, err)
    }
    defer file.Close()

    return FromReader(file, options...)
}


func (dis *Disasm) InstAt(index uint) (string, error) {
    if index < 0 || index >= uint(len(dis.romData) - 1) {
        return "", fmt.Errorf("index %d out of rom range", index)
    }

    opcode := uint16(dis.romData[index]) << 8 +
              uint16(dis.romData[index + 1])

    return ParseOpcode(opcode), nil
}


// Prefer this way to implement an iterator than using
// go routine with a channel to support `range` keyword
// since it can leak a go routine in case of early exit from iter loop
func (dis *Disasm) Iterator() iterator {
    return iterator{
        dis:   dis,

        // Since we use uint and can't cast -InstructionSize into uint,
        // so we implement the two's complement
        index: ^InstructionSize + 1,
    }
}


type iterator struct {
    dis   *Disasm
    index uint
}


func (iter *iterator) Next() bool {
    iter.index += InstructionSize

    return iter.index < uint(len(iter.dis.romData) - 1)
}


func (iter *iterator) Value() string {
    // No need to check for error since we have a check at 'Next' function
    inst, _  := iter.dis.InstAt(iter.index)

    return inst
}
