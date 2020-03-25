package godtp

import (
	"fmt"
	"testing"
	"time"
)

func TestGoDTPChan(t *testing.T) {
	waitTime := 100 * time.Millisecond

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

	// Disconnect from server
	time.Sleep(waitTime)
	err = client.Disconnect()
	assertErr(err == nil, t, err)

	// Stop server
	time.Sleep(waitTime)
	err = server.Stop()
	assertErr(err == nil, t, err)
}
