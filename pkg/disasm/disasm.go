package disasm

import (
	"io/ioutil"
	"os"
	"fmt"
	"io"
)


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
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open %v: %w", filename, err)
	}
	defer file.Close()

	return FromReader(file)
}
