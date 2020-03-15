package godtp

import "math"

const lenSize = 5

// Convert decimal to ASCII
func decToASCII(dec uint64) []byte {
	ascii := make([]byte, lenSize)
	for i := 0; i < lenSize; i++ {
		ascii[i] = uint8(dec / uint64(math.Pow(256, float64(lenSize - i - 1))))
	}
	return ascii
}

// Convert ASCII to decimal
func asciiToDec(ascii []byte) uint64 {
	var dec uint64 = 0
	for i := 0; i < lenSize; i++ {
		dec += uint64(float64(uint8(ascii[i])) * math.Pow(256, float64(lenSize - i - 1)))
	}
	return dec
}
