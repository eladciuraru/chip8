package main

import (
    "fmt"
    "flag"
    "os"
    "github.com/eladciuraru/chip8/pkg/disasm"
)


type Arguments struct {
    path   string
    base   uint
    labels bool
    format string
}


func parseArgs() Arguments {
    const (
		pathHelp   = "Path to file contains CHIP-8 compiled instructions"
		baseHelp   = "The base address of the first instruction"
		labelsHelp = "Generate labels for control flows"
		formatHelp = "Format string for each line in the disassembly, read docs for more"
		shortHelp  = " (shorthand)"
    )

    var args Arguments

    flag.StringVar(&args.path, "path", args.path, pathHelp)
    flag.StringVar(&args.path, "p", args.path, pathHelp + shortHelp)

    flag.UintVar(&args.base, "base", args.base, baseHelp)
    flag.UintVar(&args.base, "b", args.base, baseHelp + shortHelp)

    flag.BoolVar(&args.labels, "labels", args.labels, labelsHelp)
    flag.BoolVar(&args.labels, "l", args.labels, labelsHelp + shortHelp)

    flag.StringVar(&args.format, "format", args.format, formatHelp)
    flag.StringVar(&args.format, "f", args.format, formatHelp + shortHelp)

    flag.Parse()

    // Make path argument required
    if args.path == "" {
        flag.PrintDefaults()
        os.Exit(1)
    }

    return args
}


func fromArgsToOptions(args Arguments) []disasm.Option {
    var disOptions []disasm.Option
    if args.base != 0 {
        disOptions = append(disOptions, disasm.WithAddress(args.base))
    }
    if args.format == "" {
        disOptions = append(disOptions, disasm.WithFormat(args.format))
    }
    if args.labels {
        disOptions = append(disOptions, disasm.WithLabels)
    }

    return disOptions
}


func main() {
    args := parseArgs()

    dis, err := disasm.FromFile(args.path, fromArgsToOptions(args)...)
    if err != nil {
        panic(err)
    }

    // fmt.Printf("%#v\n", dis)
    // fmt.Println(dis.InstAt(0))
    // fmt.Println(dis.InstAt(2))
    // fmt.Println(dis.InstAt(3))

    for iter := dis.Iterator(); iter.Next(); {
        fmt.Println(iter.Value())
    }
}
