package disasm


func decodeArg1(opcode uint16) uint16 {
    return (opcode & 0x0F00) >> 8
}


func decodeArg2(opcode uint16) uint16 {
    return opcode & 0x00FF
}


func decodeArg3(opcode uint16) uint16 {
    return opcode & 0x0FFF
}


func decodeArgsMid(opcode uint16) (uint16, uint16) {
    return (opcode & 0x0F00) >> 8,
           (opcode & 0x00F0) >> 4
}


func decodeMajor(opcode uint16) uint16 {
    return opcode & 0xF000
}


func decodeOpcode(opcode uint16) uint16 {
    decoded := decodeMajor(opcode)

    switch decoded {
        case 0x0000: fallthrough
        case 0xE000: fallthrough
        case 0xF000: decoded = opcode & 0xF0FF
        case 0x8000: decoded = opcode & 0xF00F
    }

    return decoded
}
