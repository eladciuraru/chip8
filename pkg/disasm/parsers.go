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
}
