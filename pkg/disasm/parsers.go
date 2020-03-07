package disasm

import "fmt"

func ParseOpcode(opcode uint16) string {
    parser, ok := parsersMap[decodeOpcode(opcode)]
    if !ok {
        parser = defaultParser
    }

    return parser(opcode)
}


func defaultParser(opcode uint16) string {
    return fmt.Sprintf("%#04x", opcode)
}


type instParser func(uint16) string

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


func clsParser(opcode uint16) string {
    return "CLS"
}


func retParser(opcode uint16) string {
    return "RET"
}


func jmpParser(opcode uint16) string {
    var args string
    switch decodeMajor(opcode) {
        case 0x1000:
            args = fmt.Sprintf("%04Xh", decodeArg3(opcode))

        case 0xB000:
            args = fmt.Sprintf("V0, %04Xh", decodeArg3(opcode))
    }

    return fmt.Sprintf("JMP   %s", args)
}


func callParser(opcode uint16) string {
    return fmt.Sprintf("CALL  %04Xh", decodeArg3(opcode))
}


func seParser(opcode uint16) string {
    var args string
    switch decodeMajor(opcode) {
        case 0x3000:
            args = fmt.Sprintf("V%x, %02Xh", decodeArg1(opcode),
                               decodeArg2(opcode))

        case 0x5000:
            arg1, arg2 := decodeArgsMid(opcode)
            args = fmt.Sprintf("V%x, V%x", arg1, arg2)
    }

    return fmt.Sprintf("SE    %s", args)
}


func sneParser(opcode uint16) string {
    var args string
    switch decodeMajor(opcode) {
        case 0x4000:
            args = fmt.Sprintf("V%x, %02Xh", decodeArg1(opcode),
                               decodeArg2(opcode))

        case 0x9000:
            arg1, arg2 := decodeArgsMid(opcode)
            args = fmt.Sprintf("V%x, V%x", arg1, arg2)
    }

    return fmt.Sprintf("SNE   %s", args)
}


func ldParser(opcode uint16) string {
    var args string
    switch decodeMajor(opcode) {
        case 0x6000:
            args = fmt.Sprintf("V%x, %02Xh", decodeArg1(opcode),
                               decodeArg2(opcode))

        case 0x8000:
            arg1, arg2 := decodeArgsMid(opcode)
            args = fmt.Sprintf("V%x, V%x", arg1, arg2)

        case 0xA000:
            args = fmt.Sprintf("I, %04Xh", decodeArg3(opcode))

        case 0xF000:
            switch decodeArg2(opcode) {
                case 0x0007: args = fmt.Sprintf("V%x, DT",  decodeArg1(opcode))
                case 0x000A: args = fmt.Sprintf("V%x, K",   decodeArg1(opcode))
                case 0x0015: args = fmt.Sprintf("DT, V%x",  decodeArg1(opcode))
                case 0x0018: args = fmt.Sprintf("ST, V%x",  decodeArg1(opcode))
                case 0x0029: args = fmt.Sprintf("F, V%x",   decodeArg1(opcode))
                case 0x0033: args = fmt.Sprintf("B, V%x",   decodeArg1(opcode))
                case 0x0055: args = fmt.Sprintf("[I], V%x", decodeArg1(opcode))
                case 0x0065: args = fmt.Sprintf("V%x, [I]", decodeArg1(opcode))
            }
    }

    return fmt.Sprintf("LD    %s", args)
}


func addParser(opcode uint16) string {
    var args string
    switch decodeMajor(opcode) {
        case 0x7000:
            args = fmt.Sprintf("V%x, %02Xh", decodeArg1(opcode),
                               decodeArg2(opcode))

        // There is no need to check lower opcode part when we have only 1 option
        // because of how the map to parser func works
        case 0xF000:
            args = fmt.Sprintf("I, V%x", decodeArg1(opcode))

        case 0x8000:
            arg1, arg2 := decodeArgsMid(opcode)
            args = fmt.Sprintf("V%x, V%x", arg1, arg2)
    }

    return fmt.Sprintf("ADD   %s", args)
}


func rndParser(opcode uint16) string {
    return fmt.Sprintf("RND   V%x, %02Xh", decodeArg1(opcode), decodeArg2(opcode))
}


func drwParser(opcode uint16) string {
    arg1, arg2 := decodeArgsMid(opcode)
    arg3 := opcode & 0x000F  // The only instruction with this extra arg

    return fmt.Sprintf("DRW   V%x, V%x, %X", arg1, arg2, arg3)
}


func skpParser(opcode uint16) string {
    return fmt.Sprintf("SKP   V%x", decodeArg1(opcode))
}


func sknpParser(opcode uint16) string {
    return fmt.Sprintf("SKNP  V%x", decodeArg1(opcode))
}


func orParser(opcode uint16) string {
    arg1, arg2 := decodeArgsMid(opcode)

    return fmt.Sprintf("OR    V%x, V%x", arg1, arg2)
}


func andParser(opcode uint16) string {
    arg1, arg2 := decodeArgsMid(opcode)

    return fmt.Sprintf("AND   V%x, V%x", arg1, arg2)
}


func xorParser(opcode uint16) string {
    arg1, arg2 := decodeArgsMid(opcode)

    return fmt.Sprintf("XOR   V%x, V%x", arg1, arg2)
}


func subParser(opcode uint16) string {
    arg1, arg2 := decodeArgsMid(opcode)

    return fmt.Sprintf("SUB   V%x, V%x", arg1, arg2)
}


func shrParser(opcode uint16) string {
    return fmt.Sprintf("SHR   V%x", decodeArg1(opcode))
}


func subnParser(opcode uint16) string {
    arg1, arg2 := decodeArgsMid(opcode)

    return fmt.Sprintf("SUBN  V%x, V%x", arg1, arg2)
}


func shlParser(opcode uint16) string {
    return fmt.Sprintf("SHL   V%x", decodeArg1(opcode))
}
