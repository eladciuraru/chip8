package machine

// Why isn't there a builtin for this, seems stupid that a
// high level language like go has many things like this missing
func minUint(a, b uint) uint {
    if a <= b {
        return a
    } else {
        return b
    }
}
