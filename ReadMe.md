# CHIP8
This project is my implementation for the CHIP-8 virtual machine and developer tools,
written in Golang.

The motivation for this project was to learn a new programming language,
and my interest in writing emulators.


# Tools

### Existing Tools
* __`disasm`__ - Disassmbler of the CHIP8 instruction set

### Planned Tools
* __`casm`__ - Assembler
* __`cdbg`__ - Debugger
* __`cvm`__ - Virtual Machine

__Note__:
All the tools, will be also available as a package with an easy to use API.


# Build
```
$ build.bat <tool_name>
```
To build a specific tool, just pass the name of the tool as parameter to the build script.
When no parameter is passed, the default is to build all the tools.

__Note__:
Development was done with Go version 1.14,
and the code uses `%w` format string to wrap errors, which means that older Go version that doesn't support this will not work with the code.
