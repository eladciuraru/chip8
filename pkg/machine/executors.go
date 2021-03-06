package machine

type executor func(*Processor, uint16)


// decoded opcode to exectutor function
var executorsMap = map[uint16]executor{
    0x00E0: clsExecutor,
    0x00EE: retExecutor,
    0x1000: jmpExecutor,
    0x2000: callExecutor,
    0x3000: seExecutor,
    0x4000: sneExecutor,
    0x5000: seExecutor,
    0x6000: ldExecutor,
    0x7000: addExecutor,
    0x8000: ldExecutor,
    0x8001: orExecutor,
    0x8002: andExecutor,
    0x8003: xorExecutor,
    0x8004: addExecutor,
    0x8005: subExecutor,
    0x8006: shrExecutor,
    0x8007: subnExecutor,
    0x800E: shlExecutor,
    0x9000: sneExecutor,
    0xA000: ldExecutor,
    0xB000: jmpExecutor,
    0xC000: rndExecutor,
    0xD000: drwExecutor,
    0xE09E: skpExecutor,
    0xE0A1: sknpExecutor,
    0xF007: ldExecutor,
    // 0xF00A: ldExecutor,
    0xF015: ldExecutor,
    0xF018: ldExecutor,
    0xF01E: addExecutor,
    0xF029: ldExecutor,
    0xF033: ldExecutor,
    0xF055: ldExecutor,
    0xF065: ldExecutor,
}


func clsExecutor(cpu *Processor, opcode uint16) {
    for addr := MemoryDisplayAddr; addr < MemorySize; addr++ {
        cpu.Write(addr, 0)
    }
}


func retExecutor(cpu *Processor, opcode uint16) {
    cpu.PC = cpu.Stack[cpu.SP]

    // We want to loop around in case of over return
    cpu.SP = (cpu.SP + 1) % byte(len(cpu.Stack))
}


func jmpExecutor(cpu *Processor, opcode uint16) {
    var newPC uint16
    switch DecodeMajor(opcode) {
        case 0x1000:
            newPC = DecodeArg3(opcode)

        case 0xB000:
            newPC = uint16(cpu.V[0]) + DecodeArg3(opcode)
    }

    cpu.PC = newPC
}


func callExecutor(cpu *Processor, opcode uint16) {
    op1 := DecodeArg3(opcode)

    // We want to loop around in case of over push
    if cpu.SP == 0 {
        cpu.SP = byte(len(cpu.Stack))
    }
    cpu.SP -= 1

    cpu.Stack[cpu.SP] = cpu.PC
    cpu.PC = op1
}


func seExecutor(cpu *Processor, opcode uint16) {
    var op1, op2 uint16
    switch DecodeMajor(opcode) {
        case 0x3000:
            op1 = uint16(cpu.V[DecodeArg1(opcode)])
            op2 = DecodeArg2(opcode)

        case 0x5000:
            x, y := DecodeArgsMid(opcode)
            op1 = uint16(cpu.V[x])
            op2 = uint16(cpu.V[y])
    }

    if op1 == op2 {
        cpu.AdvancePC()
    }
}


func sneExecutor(cpu *Processor, opcode uint16) {
    var op1, op2 uint16
    switch DecodeMajor(opcode) {
        case 0x4000:
            op1 = uint16(cpu.V[DecodeArg1(opcode)])
            op2 = DecodeArg2(opcode)

        case 0x9000:
            x, y := DecodeArgsMid(opcode)
            op1 = uint16(cpu.V[x])
            op2 = uint16(cpu.V[y])
    }

    if op1 != op2 {
        cpu.AdvancePC()
    }
}


func ldExecutor(cpu *Processor, opcode uint16) {
    switch DecodeMajor(opcode) {
        case 0x6000:
            op1 := DecodeArg1(opcode)
            op2 := DecodeArg2(opcode)

            cpu.V[op1] = byte(op2)

        case 0x8000:
            op1, op2 := DecodeArgsMid(opcode)

            cpu.V[op1] = cpu.V[op2]

        case 0xA000:
            cpu.I = DecodeArg3(opcode)

        case 0xF000:
            op1 := DecodeArg1(opcode)
            switch DecodeArg2(opcode) {
                case 0x0007: cpu.V[op1] = cpu.DT
                // case 0x000A:
                //     // This opcde raised a question,
                //     // how to wait for a keyboard input without blocking
                //     // current go routine.
                //     // 3 options came to mind:
                //     // 1. polling for keypress
                //     // 2. implementing an 'interrupt' that can be activated
                //     //    from the machine API - using a channel
                //     // ^ Both require the code to not run on the main go routine
                //     // 3. 'Freeze' the state of the CPU, and each cycle will
                //     //     traverse the same code path up to here, then 'unfreeze'
                //     //     when the condition is met (pressing a key in this case)
                //     // ^ This was chosen
                //     cpu.Freeze()
                //     for i := byte(0); i < KeyF; i++ {
                //         addr := MemoryKeyboardAddr + uint16(i)
                //         if cpu.Read(addr) == KeyPressedValue {
                //             cpu.UnFreeze()
                //             cpu.V[op1] = i
                //         }
                //     }

                case 0x0015: cpu.DT = cpu.V[op1]
                case 0x0018: cpu.ST = cpu.V[op1]
                case 0x0029: cpu.I  = MemoryFontTableAddr +
                                      uint16(cpu.V[op1]) * FontSpriteSize
                case 0x0033:
                    addr  := cpu.I
                    value := cpu.V[op1]

                    cpu.Write(addr, value / 100)  // Hundreds digit
                    cpu.Write(addr + 1, value / 10 % 10)  // tens digit
                    cpu.Write(addr + 2, value % 10)  // ones digit

                case 0x0055:
                    op1 := DecodeArg1(opcode)
                    for i := uint16(0); i <= op1; i++ {
                        cpu.Write(cpu.I + i, cpu.V[i])
                    }

                case 0x0065:
                    op1 := DecodeArg1(opcode)
                    for i := uint16(0); i <= op1; i++ {
                        cpu.V[i] = cpu.Read(cpu.I + i)
                    }
            }
    }
}


func addExecutor(cpu *Processor, opcode uint16) {
    switch DecodeMajor(opcode) {
        case 0x7000:
            op1 := DecodeArg1(opcode)
            op2 := DecodeArg2(opcode)

            cpu.V[op1] += byte(op2)

        case 0x8000:
            op1, op2 := DecodeArgsMid(opcode)

            sum := uint16(cpu.V[op1]) + uint16(cpu.V[op2])
            if sum > 255 {
                cpu.V[0x0F] = 1
            } else {
                cpu.V[0x0F] = 0
            }

            cpu.V[op1] = byte(sum)

        case 0xF000:
            op1 := DecodeArg1(opcode)

            cpu.I += uint16(cpu.V[op1])
    }
}


func rndExecutor(cpu *Processor, opcode uint16) {
    op1 := DecodeArg1(opcode)
    op2 := DecodeArg2(opcode)

    cpu.V[op1] = randomByte() & byte(op2)
}


func skpExecutor(cpu *Processor, opcode uint16) {
    op1 := DecodeArg1(opcode)

    if cpu.Read(MemoryKeyboardAddr + op1) == KeyPressedValue {
        cpu.AdvancePC()
    }
}


func sknpExecutor(cpu *Processor, opcode uint16) {
    op1 := DecodeArg1(opcode)

    if cpu.Read(MemoryKeyboardAddr + op1) == KeyReleasedValue {
        cpu.AdvancePC()
    }
}


func orExecutor(cpu *Processor, opcode uint16) {
    op1, op2 := DecodeArgsMid(opcode)

    cpu.V[op1] |= cpu.V[op2]
}


func andExecutor(cpu *Processor, opcode uint16) {
    op1, op2 := DecodeArgsMid(opcode)

    cpu.V[op1] &= cpu.V[op2]
}


func xorExecutor(cpu *Processor, opcode uint16) {
    op1, op2 := DecodeArgsMid(opcode)

    cpu.V[op1] ^= cpu.V[op2]
}


func subExecutor(cpu *Processor, opcode uint16) {
    op1, op2 := DecodeArgsMid(opcode)

    if cpu.V[op1] > cpu.V[op2] {
        cpu.V[0x0F] = 1
    } else {
        cpu.V[0x0F] = 0
    }

    cpu.V[op1] -= cpu.V[op2]
}


func shrExecutor(cpu *Processor, opcode uint16) {
    op1, _ := DecodeArgsMid(opcode)

    cpu.V[0x0F]  = cpu.V[op1] & 1
    cpu.V[op1] >>= 1
}


func subnExecutor(cpu *Processor, opcode uint16) {
    op1, op2 := DecodeArgsMid(opcode)

    if cpu.V[op2] > cpu.V[op1] {
        cpu.V[0x0F] = 1
    } else {
        cpu.V[0x0F] = 0
    }

    cpu.V[op1] = cpu.V[op2] - cpu.V[op1]
}


func shlExecutor(cpu *Processor, opcode uint16) {
    op1, _ := DecodeArgsMid(opcode)

    cpu.V[0x0F]  = (cpu.V[op1] & 0x80) >> 7
    cpu.V[op1] <<= 1
}


func drwExecutor(cpu *Processor, opcode uint16) {
    const (
        collisionOn    byte = 0x01
        collisionOff   byte = 0x00
        collisionIndex byte = 0x0F
    )

    op1, op2 := DecodeArgsMid(opcode)
    op3      := opcode & 0x000F

    x, y   := uint16(cpu.V[op1]), uint16(cpu.V[op2])
    offset := y * DisplayWidth + x
    index  := offset / DisplayPixelWidth

    shrBits := (offset % DisplayPixelWidth)
    shlBits := DisplayPixelWidth - shrBits

    cpu.V[collisionIndex] = collisionOff
    for yOff := uint16(0); yOff < op3; yOff++ {
        memAddr1  := MemoryDisplayAddr + index
        memAddr2  := memAddr1 + 1
        if memAddr2 % DisplayPixelWidth == 0 {
            memAddr2 -= DisplayMemoryStride
        }

        dispData1 := cpu.Read(memAddr1)
        dispData2 := cpu.Read(memAddr2)

        spriteData  := cpu.Read(cpu.I + yOff)
        spriteData1 := spriteData >> shrBits
        spriteData2 := spriteData << shlBits

        xorData1 := dispData1 ^ spriteData1
        xorData2 := dispData2 ^ spriteData2

        // Check for collision - means we erased at least 1 pixel
        if dispData1 & ^xorData1 != 0 || dispData2 & ^xorData2 != 0 {
            cpu.V[collisionIndex] = collisionOn
        }

        cpu.Write(memAddr1, xorData1)
        cpu.Write(memAddr2, xorData2)

        index += DisplayMemoryStride
    }
}
