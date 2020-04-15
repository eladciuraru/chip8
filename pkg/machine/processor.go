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
    PC     uint16
    SP     byte
    Stack  [16]uint16  // We could store this on the RAM instead

    // The freeze mechanic enable the CPU
    freeze bool
}


func NewProcessor(bus Bus) *Processor {
    return &Processor{
        Bus: bus,
        PC:  MemoryRomAddr,
    }
}


func (cpu *Processor) Cycle() {
    cpu.HandleTimers()

    opcode := cpu.FetchOpcode()
    cpu.AdvancePC()

    cpu.Execute(opcode)
}


func (cpu *Processor) FetchOpcode() uint16 {
    opcode := uint16(cpu.Read(cpu.PC)) << 8 +
              uint16(cpu.Read(cpu.PC + 1))

    return opcode
}


func (cpu *Processor) Execute(opcode uint16) {
    executor, ok := executorsMap[DecodeOpcode(opcode)]
    if !ok {
        // TODO: Remove this when it is decided how to handle this
        // panic(fmt.Errorf("Unknown opcode found: %#04x", opcode))
        fmt.Printf("Unknown opcode found: %#04x\n", opcode)
    } else {
        // fmt.Printf("Executing instruction at: %d\n", (cpu.PC - 0x200) / 2)
        executor(cpu, opcode)
    }
}


func (cpu *Processor) AdvancePC() {
    // TODO: Handle edge case
    cpu.PC += InstructionSize
}


func (cpu *Processor) HandleTimers() {
    if cpu.DT != 0 {
        cpu.DT -= 1
    }

    if cpu.ST != 0 {
        cpu.ST -= 1
    }
}

func (cpu *Processor) restorePC() {
    // TODO: Handle edge case
    cpu.PC -= InstructionSize
}


func (cpu *Processor) Freeze() {
    cpu.freeze = true
    cpu.restorePC()
}


func (cpu *Processor) UnFreeze() {
    cpu.freeze = false
}
