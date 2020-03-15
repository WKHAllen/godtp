package godtp

import "math"

// LENSIZE defines the length of the size portion of a packet
const LENSIZE = 5

// Convert decimal to ASCII
func decToASCII(dec uint64) []byte {
	ascii := make([]byte, LENSIZE)
	for i := 0; i < LENSIZE; i++ {
		ascii[i] = uint8(dec / uint64(math.Pow(256, float64(LENSIZE - i - 1))))
	}
	return ascii
}

// Convert ASCII to decimal
func asciiToDec(ascii []byte) uint64 {
	var dec uint64 = 0
	for i := 0; i < LENSIZE; i++ {
		dec += uint64(float64(uint8(ascii[i])) * math.Pow(256, float64(LENSIZE - i - 1)))
	}
	return dec
}
