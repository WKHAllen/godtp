package godtp

const lenSize = 5

// Convert decimal to ASCII
func decToASCII(dec uint64) []byte {
	ascii := make([]byte, lenSize)
	for i := lenSize - 1; i >= 0; i-- {
		ascii[i] = uint8(dec)
		dec >>= 8
	}
	return ascii
}

// Convert ASCII to decimal
func asciiToDec(ascii []byte) uint64 {
	var dec uint64 = 0
	for i := 0; i < lenSize; i++ {
		dec <<= 8
		dec += uint64(ascii[i])
	}
	return dec
}
