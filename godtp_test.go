package godtp

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func assert(value bool, t *testing.T, err string) {
	if !value {
		t.Errorf(err)
	}
}

func assertErr(value bool, t *testing.T, err error) {
	if !value {
		t.Errorf(err.Error())
	}
}

func TestGoDTP(t *testing.T) {
	waitTime := 100 * time.Millisecond

	// Create server
	server := NewServer(nil, nil, nil, false, false)
	assert(!server.Serving(), t, "Server should not be serving")

	// Start server
	err := server.StartDefault()
	assertErr(err == nil, t, err)
	assert(server.Serving(), t, "Server should not be serving")

	// Check server address info
	host, port, err := server.GetAddr()
	assert(err == nil, t, "Error getting server address")
	assert(host == "[::]", t, "Incorrect host: " + host)
	assert(port >= 0 && port <= 65535, t, "Invalid port number" + strconv.Itoa(int(port)))
	assert(server.sock.Addr().String() == host + ":" + strconv.Itoa(int(port)), t, "Address strings don't match")
	fmt.Printf("Server running at %s:%d\n", host, port)

	// Send data
	err = server.Send([]byte{1, 3, 6})
	assertErr(err == nil, t, err)
	err = server.Send([]byte{1, 3, 6}, 0)
	assert(err.Error() == "client does not exist", t, "Send error expected")

	// Stop server
	time.Sleep(waitTime)
	err = server.Stop()
	assertErr(err == nil, t, err)
	assert(!server.Serving(), t, "Server should not be serving")
}
