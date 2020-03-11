package machine

import "crypto/rand"

// Why isn't there a builtin for this, seems stupid that a
// high level language like go has many things like this missing
func minUint16(a, b uint16) uint16 {
    if a <= b {
        return a
    } else {
        return b
    }
}


func randomByte() byte {
    buffer := make([]byte, 1)

    rand.Read(buffer)

    return buffer[0]
}
