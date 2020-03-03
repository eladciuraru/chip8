package main

import (
	"fmt"
	"github.com/eladciuraru/chip8/pkg/disasm"
)


func main() {
	dis, err := disasm.FromFile("roms/maze.bin")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", dis)
}