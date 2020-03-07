package machine

type executor func(*Processor, uint16)


// decoded opcode to exectutor function
var executorsMap = map[uint16]executor{
    // 0x00E0: clsExecutor,
    0x00EE: retExecutor,
    // 0x1000: jmpExecutor,
    0x2000: callExecutor,
    // 0x3000: seExecutor,
    // 0x4000: sneExecutor,
    // 0x5000: seExecutor,
    // 0x6000: ldExecutor,
    // 0x7000: addExecutor,
    // 0x8000: ldExecutor,
    // 0x8001: orExecutor,
    // 0x8002: andExecutor,
    // 0x8003: xorExecutor,
    // 0x8004: addExecutor,
    // 0x8005: subExecutor,
    // 0x8006: shrExecutor,
    // 0x8007: subnExecutor,
    // 0x800E: shlExecutor,
    // 0x9000: sneExecutor,
    // 0xA000: ldExecutor,
    // 0xB000: jmpExecutor,
    // 0xC000: rndExecutor,
    // 0xD000: drwExecutor,
    // 0xE09E: skpExecutor,
    // 0xE0A1: sknpExecutor,
    // 0xF007: ldExecutor,
    // 0xF00A: ldExecutor,
    // 0xF015: ldExecutor,
    // 0xF018: ldExecutor,
    // 0xF01E: addExecutor,
    // 0xF029: ldExecutor,
    // 0xF033: ldExecutor,
    // 0xF055: ldExecutor,
    // 0xF065: ldExecutor,
}


func clsExecutor(cpu *Processor, opcode uint16) {

}


func retExecutor(cpu *Processor, opcode uint16) {
    cpu.PC = cpu.Stack[cpu.SP] 

    // We want to loop around in case of over return
    cpu.SP = (cpu.SP + 1) % byte(len(cpu.Stack))
}


func jmpExecutor(cpu *Processor, opcode uint16) {
    
}


func callExecutor(cpu *Processor, opcode uint16) {
    // We want to loop around in case of over push
    if cpu.SP == 0 {
        cpu.SP = byte(len(cpu.Stack))
    }
    cpu.SP -= 1

    cpu.Stack[cpu.SP] = cpu.PC
    cpu.PC = DecodeArg3(opcode)
}