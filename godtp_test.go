package godtp

import (
	"strings"
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

	assert(server.GetNetwork() == "tcp", t, "Network should be 'tcp'")
	assert(strings.HasPrefix(server.GetAddr(), "[::]:"), t, "Address should be '[::]:XXXXX'")
	assert(server.GetHost() == "[::]", t, "Host should be '[::]'")
	assert(server.GetPort() >= 0 && server.GetPort() <= 65535, t, "Invalid port number")

	time.Sleep(waitTime)
	err = server.Stop()
	assertErr(err == nil, t, err)
	assert(!server.Serving(), t, "Server should not be serving")
}
