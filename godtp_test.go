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

func onRecvServer(clientID uint, data []byte) {
	fmt.Printf("Data received from client #%d: %s\n", clientID, data)
}

func onConnectServer(clientID uint) {
	fmt.Printf("Client #%d connected\n", clientID)
}

func onDisconnectServer(clientID uint) {
	fmt.Printf("Client #%d disconnected\n", clientID)
}

func onRecvClient(data []byte) {
	fmt.Printf("Data received from server: %s\n", data)
}

func onDisconnectedClient() {
	fmt.Printf("Disconnected from server\n")
}

func TestGoDTP(t *testing.T) {
	waitTime := 100 * time.Millisecond

	// Create server
	server := NewServer(onRecvServer, onConnectServer, onDisconnectServer, false, false)
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

	// Create client
	client := NewClient(onRecvClient, onDisconnectedClient, false, false)
	assert(!client.Connected(), t, "Client should not be connected")

	// Connect to server
	err = client.Connect(host, port)
	assertErr(err == nil, t, err)

	// Check that addresses match
	host1, port1, err := client.GetAddr()
	host2, port2, err := server.GetClientAddr(0)
	assert(host1 == host2, t, "Client hosts do not match")
	assert(port1 == port2, t, "Client ports do not match")
	_, port3, err := client.GetServerAddr()
	_, port4, err := server.GetAddr()
	assert(port3 == port4, t, "Server ports do not match")

	// Disconnect from server
	time.Sleep(waitTime)
	err = client.Disconnect()
	assertErr(err == nil, t, err)

	// Stop server
	time.Sleep(waitTime)
	err = server.Stop()
	assertErr(err == nil, t, err)
	assert(!server.Serving(), t, "Server should not be serving")
}
