package disasm

func ParseOpcode(opcode uint16) *Instruction {
    parser, ok := parsersMap[decodeOpcode(opcode)]
    if !ok {
        parser = defaultParser
    }

    return parser(opcode)
}


func defaultParser(opcode uint16) *Instruction {
    return newInstruction(opcode, "")
}


type instParser func(uint16) *Instruction


// decoded opcode to instParser function
var parsersMap = map[uint16]instParser{
    0x00E0: clsParser,
    0x00EE: retParser,
    0x1000: jmpParser,
    0x2000: callParser,
    0x3000: seParser,
    0x4000: sneParser,
    0x5000: seParser,
    0x6000: ldParser,
    0x7000: addParser,
    0x8000: ldParser,
    0x8001: orParser,
    0x8002: andParser,
    0x8003: xorParser,
    0x8004: addParser,
    0x8005: subParser,
    0x8006: shrParser,
    0x8007: subnParser,
    0x800E: shlParser,
    0x9000: sneParser,
    0xA000: ldParser,
    0xB000: jmpParser,
    0xC000: rndParser,
    0xD000: drwParser,
    0xE09E: skpParser,
    0xE0A1: sknpParser,
    0xF007: ldParser,
    0xF00A: ldParser,
    0xF015: ldParser,
    0xF018: ldParser,
    0xF01E: addParser,
    0xF029: ldParser,
    0xF033: ldParser,
    0xF055: ldParser,
    0xF065: ldParser,
}


func clsParser(opcode uint16) *Instruction {
    return newInstruction(opcode, MnemonicCLS)
}


func retParser(opcode uint16) *Instruction {
    return newInstruction(opcode, MnemonicRET)
}


func jmpParser(opcode uint16) *Instruction {
    inst := newInstruction(opcode, MnemonicJMP)
    switch decodeMajor(opcode) {
        case 0x1000:
            inst.Operand1 = Immediate(decodeArg3(opcode))

        case 0xB000:
            inst.Operand1 = RegisterV0
            inst.Operand2 = Immediate(decodeArg3(opcode))
    }

    return inst
}


func callParser(opcode uint16) *Instruction {
    inst := newInstruction(opcode, MnemonicCALL)
    inst.Operand1 = Immediate(decodeArg3(opcode))

    return inst
}


func seParser(opcode uint16) *Instruction {
    inst := newInstruction(opcode, MnemonicSE)
    switch decodeMajor(opcode) {
        case 0x3000:
            inst.Operand1 = vRegisters[decodeArg1(opcode)]
            inst.Operand2 = Immediate(decodeArg2(opcode))

        case 0x5000:
            arg1, arg2 := decodeArgsMid(opcode)
            inst.Operand1 = vRegisters[arg1]
            inst.Operand2 = vRegisters[arg2]
    }

    return inst
}


func sneParser(opcode uint16) *Instruction {
    inst := newInstruction(opcode, MnemonicSNE)
    switch decodeMajor(opcode) {
        case 0x4000:
            inst.Operand1 = vRegisters[decodeArg1(opcode)]
            inst.Operand2 = Immediate(decodeArg2(opcode))

        case 0x9000:
            arg1, arg2 := decodeArgsMid(opcode)
            inst.Operand1 = vRegisters[arg1]
            inst.Operand2 = vRegisters[arg2]
    }

    return inst
}


func ldParser(opcode uint16) *Instruction {
    inst := newInstruction(opcode, MnemonicLD)
    switch decodeMajor(opcode) {
        case 0x6000:
            inst.Operand1 = vRegisters[decodeArg1(opcode)]
            inst.Operand2 = Immediate(decodeArg2(opcode))

        case 0x8000:
            arg1, arg2 := decodeArgsMid(opcode)
            inst.Operand1 = vRegisters[arg1]
            inst.Operand2 = vRegisters[arg2]

        case 0xA000:
            inst.Operand1 = RegisterI
            inst.Operand2 = Immediate(decodeArg3(opcode))

        case 0xF000:
            switch decodeArg2(opcode) {
                case 0x0007:
                    inst.Operand1 = vRegisters[decodeArg1(opcode)]
                    inst.Operand2 = RegisterDT

                case 0x000A:
                    inst.Operand1 = vRegisters[decodeArg1(opcode)]
                    inst.Operand2 = RegisterKeyPress

                case 0x0015:
                    inst.Operand1 = RegisterDT
                    inst.Operand2 = vRegisters[decodeArg1(opcode)]

                case 0x0018:
                    inst.Operand1 = RegisterST
                    inst.Operand2 = vRegisters[decodeArg1(opcode)]

                case 0x0029:
                    inst.Operand1 = RegisterFont
                    inst.Operand2 = vRegisters[decodeArg1(opcode)]

                case 0x0033:
                    inst.Operand1 = RegisterBCD
                    inst.Operand2 = vRegisters[decodeArg1(opcode)]

                case 0x0055:
                    inst.Operand1 = Pointer(RegisterI)
                    inst.Operand2 = vRegisters[decodeArg1(opcode)]

                case 0x0065:
                    inst.Operand1 = vRegisters[decodeArg1(opcode)]
                    inst.Operand2 = Pointer(RegisterI)
            }
    }

    return inst
}


func addParser(opcode uint16) *Instruction {
    inst := newInstruction(opcode, MnemonicADD)
    switch decodeMajor(opcode) {
        case 0x7000:
            inst.Operand1 = vRegisters[decodeArg1(opcode)]
            inst.Operand2 = Immediate(decodeArg2(opcode))

        // There is no need to check lower opcode part when we have only 1 option
        // because of how the map to parser func works
        case 0xF000:
            inst.Operand1 = RegisterI
            inst.Operand2 = vRegisters[decodeArg1(opcode)]

        case 0x8000:
            arg1, arg2 := decodeArgsMid(opcode)
            inst.Operand1 = vRegisters[arg1]
            inst.Operand2 = vRegisters[arg2]
    }

    return inst
}


func rndParser(opcode uint16) *Instruction {
    inst := newInstruction(opcode, MnemonicRND)
    inst.Operand1 = vRegisters[decodeArg1(opcode)]
    inst.Operand2 = Immediate(decodeArg2(opcode))

    return inst
}


func drwParser(opcode uint16) *Instruction {
    inst := newInstruction(opcode, MnemonicDRW)

    arg1, arg2   := decodeArgsMid(opcode)
    inst.Operand1 = vRegisters[arg1]
    inst.Operand2 = vRegisters[arg2]
    inst.Operand3 = Immediate(opcode & 0x000F)  // The only instruction with this extra operand

    return inst
}


func skpParser(opcode uint16) *Instruction {
    inst := newInstruction(opcode, MnemonicSKP)
    inst.Operand1 = vRegisters[decodeArg1(opcode)]

    return inst
}


func sknpParser(opcode uint16) *Instruction {
    inst := newInstruction(opcode, MnemonicSKNP)
    inst.Operand1 = vRegisters[decodeArg1(opcode)]

    return inst
}


func orParser(opcode uint16) *Instruction {
    inst := newInstruction(opcode, MnemonicOR)

    arg1, arg2   := decodeArgsMid(opcode)
    inst.Operand1 = vRegisters[arg1]
    inst.Operand2 = vRegisters[arg2]

    return inst
}


func andParser(opcode uint16) *Instruction {
    inst := newInstruction(opcode, MnemonicAND)

    arg1, arg2   := decodeArgsMid(opcode)
    inst.Operand1 = vRegisters[arg1]
    inst.Operand2 = vRegisters[arg2]

    return inst
}


func xorParser(opcode uint16) *Instruction {
    inst := newInstruction(opcode, MnemonicXOR)

    arg1, arg2   := decodeArgsMid(opcode)
    inst.Operand1 = vRegisters[arg1]
    inst.Operand2 = vRegisters[arg2]

    return inst
}


func subParser(opcode uint16) *Instruction {
    inst := newInstruction(opcode, MnemonicSUB)

    arg1, arg2   := decodeArgsMid(opcode)
    inst.Operand1 = vRegisters[arg1]
    inst.Operand2 = vRegisters[arg2]

    return inst
}


func shrParser(opcode uint16) *Instruction {
    inst := newInstruction(opcode, MnemonicSHR)
    inst.Operand1 = vRegisters[decodeArg1(opcode)]

    return inst
}


func subnParser(opcode uint16) *Instruction {
    inst := newInstruction(opcode, MnemonicSUBN)

    arg1, arg2   := decodeArgsMid(opcode)
    inst.Operand1 = vRegisters[arg1]
    inst.Operand2 = vRegisters[arg2]

    return inst
}


func shlParser(opcode uint16) *Instruction {
    inst := newInstruction(opcode, MnemonicSHL)
    inst.Operand1 = vRegisters[decodeArg1(opcode)]

    return inst
}
