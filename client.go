package godtp

import (
	"net"
	"sync"
)

type onRecvFuncClient func([]byte)
type onDisconnectedFuncClient func()

// Client defines the socket client type
type Client struct {
	onRecv onRecvFuncClient
	onDisconnected onDisconnectedFuncClient
	blocking bool
	eventBlocking bool
	connected bool
	sock net.Conn
	wg sync.WaitGroup
}

// NewClient creates a new socket client object
func NewClient(onRecv onRecvFuncClient,
			   onDisconnected onDisconnectedFuncClient,
			   blocking bool, eventBlocking bool) *Client {
	return &Client{
		onRecv: onRecv,
		onDisconnected: onDisconnected,
		blocking: blocking,
		eventBlocking: eventBlocking,
		connected: false,
	}
}

// NewClientDefault creates a new socket client object with blocking and
// eventBlocking set to false
func NewClientDefault(onRecv onRecvFuncClient,
					  onDisconnected onDisconnectedFuncClient) *Client {
	return &Client{
		onRecv: onRecv,
		onDisconnected: onDisconnected,
		blocking: false,
		eventBlocking: false,
		connected: false,
	}
}

// Connected returns a boolean value representing whether or not the client is
// connected to a server
func (client *Client) Connected() bool {
	return client.connected
}
