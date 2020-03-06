package machine

type Processor struct {
    bus Bus

    // Registers
    V  [16]byte  // General registers
    I  uint16    // Used for addressing, only 12 bits are used
    DT byte
    ST byte

    // Pseduo registers
    PC    uint16
    SP    byte
    Stack [16]uint16  // We could store this on the RAM instead
}


func NewProcessor(bus Bus, startAddr uint16) *Processor {
    return &Processor{
        bus: bus,
        PC:  startAddr,
    }
}
