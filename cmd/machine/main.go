package main

import (
	"fmt"

	"github.com/eladciuraru/chip8/pkg/machine"
)


func main() {
    vm, err := machine.FromFile("roms/maze.bin")
    if err != nil {
        panic(err)
    }

    fmt.Println(vm)
}
