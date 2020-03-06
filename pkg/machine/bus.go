package machine

// This type might seems a little bit silly,
// but the reason is I need a way to communicate between
// the CPU and the RAM, and while this is over simplification
// of that task, it is still closer to how real hardware works
// rather than just letting the CPU have a raw access to the RAM
type Bus struct {
    data uint16
}


func (bus *Bus) Send(data uint16) {
    bus.data = data
}


func (bus *Bus) Read() uint16 {
    return bus.data
}
