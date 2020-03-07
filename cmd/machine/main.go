package main

import (
	"github.com/eladciuraru/chip8/pkg/machine"
)

func main() {
    vm, err := machine.FromFile("roms/maze.bin")
    if err != nil {
        panic(err)
    }

    vm.Start()
}
