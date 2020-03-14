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

func TestServer(t *testing.T) {
	waitTime := 100 * time.Millisecond

	server := NewServer(nil, nil, nil, false, false)
	assert(!server.Serving(), t, "Server should not be serving")

	err := server.StartDefault()
	assertErr(err == nil, t, err)
	assert(server.Serving(), t, "Server should not be serving")

	host, port, err := server.GetAddr()
	assert(err == nil, t, "Error getting server address")
	assert(host == "[::]", t, "Incorrect host: " + host)
	assert(port >= 0 && port <= 65535, t, "Invalid port number" + strconv.Itoa(int(port)))
	assert(server.sock.Addr().String() == host + ":" + strconv.Itoa(int(port)), t, "Address strings don't match")
	fmt.Printf("Server running at %s:%d\n", host, port)

	time.Sleep(waitTime)
	err = server.Stop()
	assertErr(err == nil, t, err)
	assert(!server.Serving(), t, "Server should not be serving")
}
