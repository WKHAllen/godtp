package godtp

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

const waitTime = 100 * time.Millisecond

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

func TestMain(t *testing.T) {
	// Create server
	server := NewServer(onRecvServer, onConnectServer, onDisconnectServer, false, false)
	assert(!server.Serving(), t, "Server should not be serving")

	// Start server
	err := server.StartDefault()
	assertErr(err == nil, t, err)
	assert(server.Serving(), t, "Server should not be serving")

	// Check server address info
	host, port, err := server.GetAddr()
	assertErr(err == nil, t, err)
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

	// Send data
	time.Sleep(waitTime)
	server.Send([]byte("Hello, client #0!"))
	time.Sleep(waitTime)
	client.Send([]byte("Hello, server!"))

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

func TestChan(t *testing.T) {
	// Create and start server
	server := NewServerChan()
	connectChan, err := server.StartDefault()
	assertErr(err == nil, t, err)

	// Get server address
	host, port, err := server.GetAddr()
	assertErr(err == nil, t, err)

	// Create client and connect to server
	client := NewClientChan()
	clientSendChan, clientRecvChan, err := client.Connect(host, port)
	assertErr(err == nil, t, err)

	// Get server send and receive channels
	time.Sleep(waitTime)
	serverSendChan := <-connectChan
	serverRecvChan := <-connectChan

	// Send data from server to client
	time.Sleep(waitTime)
	serverSendChan <- []byte("Hello, client!")
	time.Sleep(waitTime)
	clientRecv := <-clientRecvChan
	fmt.Printf("Data received from server: %s\n", clientRecv)
	assert(string(clientRecv) == "Hello, client!", t, "Unexpected data received from server")

	// Send data from client to server
	time.Sleep(waitTime)
	clientSendChan <- []byte("Hello, server!")
	time.Sleep(waitTime)
	serverRecv := <-serverRecvChan
	fmt.Printf("Data received from client: %s\n", serverRecv)
	assert(string(serverRecv) == "Hello, server!", t, "Unexpected data received from client")

	// Test channel close
	go func() {
		for range clientRecvChan {}
		fmt.Println("Disconnected from server")
	}()
	go func() {
		for range serverRecvChan {}
		fmt.Println("Client disconnected")
	}()

	// Disconnect from server
	time.Sleep(waitTime)
	err = client.Disconnect()
	assertErr(err == nil, t, err)

	// Stop server
	time.Sleep(waitTime)
	err = server.Stop()
	assertErr(err == nil, t, err)
}

func TestEncode(t *testing.T) {
	// Create and start server
	server := NewServerChan()
	connectChan, err := server.StartDefault()
	assertErr(err == nil, t, err)

	// Get server address
	host, port, err := server.GetAddr()
	assertErr(err == nil, t, err)

	// Create client and connect to server
	client := NewClientChan()
	clientSendChan, clientRecvChan, err := client.Connect(host, port)
	assertErr(err == nil, t, err)

	// Get server send and receive channels
	time.Sleep(waitTime)
	serverSendChan := <-connectChan
	serverRecvChan := <-connectChan

	// Send data from server to client
	time.Sleep(waitTime)
	enc1, err := Encode(123)
	assertErr(err == nil, t, err)
	serverSendChan <- enc1
	time.Sleep(waitTime)
	clientRecv := <-clientRecvChan
	var dec1 int
	Decode(clientRecv, &dec1)
	fmt.Printf("Data received from server: %d\n", dec1)
	assert(dec1 == 123, t, "Unexpected data received from server")

	// Send data from client to server
	time.Sleep(waitTime)
	enc2, err := Encode("Test 123")
	assertErr(err == nil, t, err)
	clientSendChan <- enc2
	time.Sleep(waitTime)
	serverRecv := <-serverRecvChan
	var dec2 string
	Decode(serverRecv, &dec2)
	fmt.Printf("Data received from client: %s\n", dec2)
	assert(dec2 == "Test 123", t, "Unexpected data received from client")

	// Disconnect from server
	time.Sleep(waitTime)
	err = client.Disconnect()
	assertErr(err == nil, t, err)

	// Stop server
	time.Sleep(waitTime)
	err = server.Stop()
	assertErr(err == nil, t, err)
}
