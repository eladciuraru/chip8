package disasm

import (
	"io/ioutil"
	"os"
	"fmt"
	"io"
)


const InstructionSize uint = 0x02


type Disasm struct {
	romData []byte
}


func New(rom []byte) *Disasm {
	return &Disasm{
		romData: rom,
	}
}


func FromReader(reader io.Reader) (*Disasm, error) {
	romData, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read all: %w", err)
	}

	return New(romData), nil
}


func FromFile(filename string) (*Disasm, error) {
	// Maybe it is better to just use ioutil.ReadFile instead of reusing
	// `FromReader` function
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open %v: %w", filename, err)
	}
	defer file.Close()

	return FromReader(file)
}


func (dis *Disasm) InstAt(index uint) (string, error) {
	if index < 0 || index >= uint(len(dis.romData) - 1) {
		return "", fmt.Errorf("index %d out of rom range", index)
	}

	opcode := uint16(dis.romData[index]) << 8 +
			  uint16(dis.romData[index + 1])

	return fmt.Sprintf("%#04x", opcode), nil
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
