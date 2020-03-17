package main

import (
	"syscall"
	"unsafe"
)

type Kernel32DLL struct {
    _GetModuleHandleW *syscall.Proc
}


var kernel32 = NewKernel32() 


func NewKernel32() *Kernel32DLL {
    dll := syscall.MustLoadDLL("kernel32")

    return &Kernel32DLL{
        _GetModuleHandleW: dll.MustFindProc("GetModuleHandleW"),
    }
}


func (k32 *Kernel32DLL) GetModuleHandleW(lpModuleName *uint16) syscall.Handle {
    ret, _, _ := k32._GetModuleHandleW.Call(uintptr(unsafe.Pointer(lpModuleName)))

    return syscall.Handle(ret)
}
