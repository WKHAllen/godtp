package godtp

import (
	"net"
	"sync"
)

type onRecvFunc func(int, []byte)
type onConnectFunc func(int)
type onDisconnectFunc func(int)

// Server defines the socket server type
type Server struct {
	onRecv onRecvFunc
	onConnect onConnectFunc
	onDisconnect onDisconnectFunc
	blocking bool
	eventBlocking bool
	serving bool
	sock net.Listener
	socks []net.Conn
	wg sync.WaitGroup
}

// NewServer creates a new socket server object
func NewServer(onRecv onRecvFunc, onConnect onConnectFunc, onDisconnect onDisconnectFunc,
			   blocking bool, eventBlocking bool) *Server {
	return &Server{
		onRecv: onRecv,
		onConnect: onConnect,
		onDisconnect: onDisconnect,
		blocking: blocking,
		eventBlocking: eventBlocking,
		serving: false,
	}
}
