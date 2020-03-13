package godtp

import (
	"fmt"
	"strconv"
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
		socks: make([]net.Conn, 0),
	}
}

// NewServerDefault creates a new socket server object with blocking and eventBlocking set to false
func NewServerDefault(onRecv onRecvFunc, onConnect onConnectFunc, onDisconnect onDisconnectFunc) *Server {
	return NewServer(onRecv, onConnect, onDisconnect, false, false)
}

// Start a server
func (server *Server) Start(host string, port uint16) error {
	if server.serving {
		return fmt.Errorf("already serving")
	}

	address := host + ":" + strconv.Itoa(int(port))
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	server.sock = ln

	server.serving = true
	if server.blocking {
		server.serve()
	} else {
		server.wg.Add(1)
		go server.serve()
	}

	return nil
}

// StartDefaultHost starts the server at the default host address
func (server *Server) StartDefaultHost(port uint16) error {
	return server.Start("0.0.0.0", port)
}

// StartDefaultPort starts the server on the default port
func (server *Server) StartDefaultPort(host string) error {
	return server.Start(host, 0)
}

// StartDefault starts the server on 0.0.0.0:0
func (server *Server) StartDefault() error {
	return server.Start("0.0.0.0", 0)
}

// Serve clients
func (server *Server) serve() {
	
}
