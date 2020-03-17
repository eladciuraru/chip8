// +build windows

package main

func main() {
    window := newWindow("test", 200, 200)
    window.MessageLoop()
    // vm, err := machine.FromFile("roms/maze.bin")
    // if err != nil {
    //     panic(err)
    // }

    // vm.SetKeyState(machine.KeyA, true)
}


func Assert(cond bool, err error) {
	if cond == false {
		panic(err)
	}
}
