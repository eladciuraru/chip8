// +build windows

package main

import "runtime"

func main() {
    // This is required, since handling WM_PAINT message requires
    // the same thread that registered the WndProc
    runtime.LockOSThread()
    
    window := NewWindow("test", 200, 200)
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
