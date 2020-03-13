package machine

import (
	"fmt"
)

type Processor struct {
    Bus

    // Registers
    V  [16]byte  // General registers
    I  uint16    // Used for addressing, only 12 bits are used
    DT byte
    ST byte

    // Pseduo registers
    PC    uint16
    SP    byte
    Stack [16]uint16  // We could store this on the RAM instead
}


func NewProcessor(bus Bus) *Processor {
    return &Processor{
        Bus: bus,
        PC:  MemoryRomAddr,
    }
}


func (cpu *Processor) Cycle() {
    opcode := cpu.FetchOpcode()

    cpu.Execute(opcode)
}


func (cpu *Processor) FetchOpcode() uint16 {
    opcode := uint16(cpu.Read(cpu.PC)) << 8 +
              uint16(cpu.Read(cpu.PC + 1))
    cpu.AdvancePC()

    return opcode
}


func (cpu *Processor) Execute(opcode uint16) {
    executor, ok := executorsMap[DecodeOpcode(opcode)]
    if !ok {
        // TODO: Remove this when it is decided how to handle this
        panic(fmt.Errorf("Unknown opcode found"))
    }

    executor(cpu, opcode)
}


func (cpu *Processor) AdvancePC() {
    // TODO: Handle edge case
    cpu.PC += InstructionSize
}
