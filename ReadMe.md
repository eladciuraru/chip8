# CHIP8
This project is my implementation for the CHIP-8 virtual machine and developer tools,
written in Golang.

The motivation for this project was to learn a new programming language,
and my interest in writing emulators.


# Tools

### Existing Tools
* __`disasm`__ - Disassmbler of the CHIP8 instruction set
* __`machine`__ - Virtual Machine

### Planned Tools
* __`casm`__ - Assembler
* __`cdbg`__ - Debugger

__Notes__:
* All the tools, will be also available as a package with an easy to use API.
* virtual machine command only supported on windows


# Build
```
$ build.bat <tool_name>
```
To build a specific tool, just pass the name of the tool as parameter to the build script.
When no parameter is passed, the default is to build all the tools.

__Note__:
Development was done with Go version 1.14,
and the code uses `%w` format string to wrap errors, which means that older Go version that doesn't support this will not work with the code.


# Emulator Design Notes
I am trying to go for a more realistic emulator implementation,
rather than an implementation for an interpreted language.
Which means the design of the emulator code will be divided into
multiple components, that will communicate with each other through the bus interface.

__note__: The bus implementation will be simple, and will not be implemented
with cycles in mind, just a simple API of passing data based of a given address.

The following is the memory map of the system:
```
0x0000 - 0x01FF : Reserved for interpreter
0x0200 - 0x0E9F : Program / Data space
0x0EA0 - 0x0EFF : Stack space
0x0ED0 - 0x0EEF : Work area
0x0EF0 - 0x0EFF : Variable space (V registers)
0x0F00 - 0x0FFF : Display buffer
```
With a couple of changes in mind:
* Variable space will be changed to keyboard state space
* Work area will not be used at all
* Since the interpreter will live outside this memory space,
  the memory reserved for the interpreter will be used to save
  some sprites data, like fonts
