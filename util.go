package godtp

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// The length of the size portion of a message
const lenSize = 5

// The buffer size of each channel
const channelBufferSize = 100

// Encode an object, so it can be sent through a socket
func encodeObject[T any](object T) ([]byte, error) {
	return json.Marshal(object)
}

// Decode an object coming from a socket
func decodeObject[T any](byteString []byte) (T, error) {
	var object T
	err := json.Unmarshal(byteString, &object)
	return object, err
}

// Encode the size portion of a message
func encodeMessageSize(size uint64) []byte {
	encodedSize := make([]byte, lenSize)

	for i := lenSize - 1; i >= 0; i-- {
		encodedSize[i] = uint8(size)
		size >>= 8
	}

	return encodedSize
}

// Decode the size portion of a message
func decodeMessageSize(encodedSize []byte) uint64 {
	var size uint64 = 0

	for i := 0; i < lenSize; i++ {
		size <<= 8
		size += uint64(encodedSize[i])
	}

	return size
}

// Parse an address
func parseAddr(addr string) (string, uint16, error) {
	index := strings.LastIndex(addr, ":")
	if index > -1 {
		port, err := strconv.Atoi(addr[index+1:])
		if err == nil {
			return addr[:index], uint16(port), nil
		}
		return "", 0, fmt.Errorf("port conversion failed")
	}
	return "", 0, fmt.Errorf("no port found")
}
