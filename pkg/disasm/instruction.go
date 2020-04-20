package disasm

import (
	"fmt"
	"strings"
)

type Instruction struct {
    // Metadata
    Address uint
    Opcode  uint16

    Mnemonic string
    Operand1 Operand
    Operand2 Operand
    Operand3 Operand
}


type Operand interface {
    operand()
    String() string
}

type Immediate uint16
type Register  string
type Pointer   Register

func (imm Immediate) operand() { }
func (reg Register)  operand() { }
func (ptr Pointer)   operand() { }

func (imm Immediate) String() string { return fmt.Sprintf("%04Xh", uint16(imm)) }
func (reg Register)  String() string { return string(reg) }
func (ptr Pointer)   String() string { return fmt.Sprintf("[%s]", string(ptr)) }


func newInstruction(opcode uint16, mnemonic string) *Instruction {
    return &Instruction{
        Opcode:   opcode,
        Mnemonic: mnemonic,
    }
}


func (inst *Instruction) String() string {
    str := fmt.Sprintf("%04Xh", inst.Address)
    if inst.Mnemonic == "" {
        return fmt.Sprintf("%s - dw    %04Xh", str, inst.Opcode)
    }

    str += fmt.Sprintf(" - %-5s", inst.Mnemonic)

    if inst.Operand1 != nil {
        str += fmt.Sprintf(" %s", inst.Operand1)
    }

    if inst.Operand2 != nil {
        spaces := strings.Repeat(" ", 5 - len(inst.Operand1.String()))
        str += fmt.Sprintf(",%s%s", spaces, inst.Operand2)
    }

    if inst.Operand3 != nil {
        str += fmt.Sprintf(",    %s", inst.Operand3)
    }

    return str
}


const (
    // Mnemonics
    MnemonicCLS  string = "cls"
    MnemonicRET  string = "ret"
    MnemonicJMP  string = "jmp"
    MnemonicCALL string = "call"
    MnemonicSE   string = "se"
    MnemonicSNE  string = "sne"
    MnemonicLD   string = "ld"
    MnemonicADD  string = "add"
    MnemonicOR   string = "or"
    MnemonicAND  string = "and"
    MnemonicXOR  string = "xor"
    MnemonicSUB  string = "sub"
    MnemonicSHR  string = "shr"
    MnemonicSUBN string = "subn"
    MnemonicSHL  string = "shl"
    MnemonicRND  string = "rnd"
    MnemonicDRW  string = "drw"
    MnemonicSKP  string = "skp"
    MnemonicSKNP string = "skpn"

    // Registers
    RegisterV0  Register = "V0"
    RegisterV1  Register = "V1"
    RegisterV2  Register = "V2"
    RegisterV3  Register = "V3"
    RegisterV4  Register = "V4"
    RegisterV5  Register = "V5"
    RegisterV6  Register = "V6"
    RegisterV7  Register = "V7"
    RegisterV8  Register = "V8"
    RegisterV9  Register = "V9"
    RegisterVa  Register = "Va"
    RegisterVb  Register = "Vb"
    RegisterVc  Register = "Vc"
    RegisterVd  Register = "Vd"
    RegisterVe  Register = "Ve"
    RegisterVf  Register = "Vf"
    RegisterI   Register = "I"
    RegisterDT  Register = "DT"
    RegisterST  Register = "ST"
    // Not real registers, but has special print char
    RegisterKeyPress Register = "K"
    RegisterFont     Register = "F"
    RegisterBCD      Register = "B"
)


var vRegisters = []Register{
    RegisterV0,
    RegisterV1,
    RegisterV2,
    RegisterV3,
    RegisterV4,
    RegisterV5,
    RegisterV6,
    RegisterV7,
    RegisterV8,
    RegisterV9,
    RegisterVa,
    RegisterVb,
    RegisterVc,
    RegisterVd,
    RegisterVe,
    RegisterVf,
}
