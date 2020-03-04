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


func decodeOpcode(opcode uint16) uint16 {
	return opcode
}
