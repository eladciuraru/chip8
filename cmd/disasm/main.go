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
	fmt.Println(dis.InstAt(0))
	fmt.Println(dis.InstAt(2))
	fmt.Println(dis.InstAt(3))

	for iter := dis.Iterator(); iter.Next(); {
		fmt.Println(iter.Value())
	}
}
