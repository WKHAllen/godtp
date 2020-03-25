package godtp

import (
	"fmt"
	"testing"
	"time"
)

func TestGoDTPChan(t *testing.T) {
	waitTime := 100 * time.Millisecond

	server := NewServerChanDefault()
	connectChan, err := server.StartDefault()
	assertErr(err == nil, t, err)

	host, port, err := server.GetAddr()
	assertErr(err == nil, t, err)

	client := NewClientDefault(onRecvClient, onDisconnectedClient)
	err = client.Connect(host, port)
	assertErr(err == nil, t, err)

	time.Sleep(waitTime)
	clientChans := <-connectChan

	time.Sleep(waitTime)
	clientChans.SendChan <- []byte("Hello, client!")
	time.Sleep(waitTime)
	client.Send([]byte("Hello, server!"))
	time.Sleep(waitTime)
	recvData := <-clientChans.RecvChan
	fmt.Printf("Data received from client: %s\n", recvData)
	assert(string(recvData) == "Hello, server!", t, "Unexpected data received from client")

	time.Sleep(waitTime)
	err = client.Disconnect()
	assertErr(err == nil, t, err)

	time.Sleep(waitTime)
	err = server.Stop()
	assertErr(err == nil, t, err)
}
