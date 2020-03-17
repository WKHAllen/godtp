package godtp

import (
	"fmt"
	"net"
	"strconv"
	"sync"
)

type onRecvFuncServer func(uint, []byte)
type onConnectFuncServer func(uint)
type onDisconnectFuncServer func(uint)

// Server defines the socket server type
type Server struct {
	onRecv onRecvFuncServer
	onConnect onConnectFuncServer
	onDisconnect onDisconnectFuncServer
	blocking bool
	eventBlocking bool
	serving bool
	sock net.Listener
	clients map[uint]net.Conn
	wg sync.WaitGroup
	nextClientID uint
}

// NewServer creates a new socket server object
func NewServer(onRecv onRecvFuncServer,
			   onConnect onConnectFuncServer,
			   onDisconnect onDisconnectFuncServer,
			   blocking bool, eventBlocking bool) *Server {
	return &Server{
		onRecv: onRecv,
		onConnect: onConnect,
		onDisconnect: onDisconnect,
		blocking: blocking,
		eventBlocking: eventBlocking,
		serving: false,
		clients: make(map[uint]net.Conn),
		nextClientID: 0,
	}
}

// NewServerDefault creates a new socket server object with blocking and
// eventBlocking set to false
func NewServerDefault(onRecv onRecvFuncServer,
					  onConnect onConnectFuncServer,
					  onDisconnect onDisconnectFuncServer) *Server {
	return NewServer(onRecv, onConnect, onDisconnect, false, false)
}

// Start the server
func (server *Server) Start(host string, port uint16) error {
	if server.serving {
		return fmt.Errorf("server is already serving")
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

// Stop the server
func (server *Server) Stop() error {
	if !server.serving {
		return fmt.Errorf("server is not serving")
	}

	server.serving = false

	for _, client := range server.clients {
		err := client.Close()
		if err != nil {
			return err
		}
	}
	err := server.sock.Close()
	if err != nil {
		return err
	}

	if !server.blocking {
		server.wg.Wait()
	}

	return nil
}

// Send data to clients
func (server *Server) Send(data []byte, clientIDs ...uint) error {
	if !server.serving {
		return fmt.Errorf("server is not serving")
	}

	size := decToASCII(uint64(len(data)))
	buffer := append(size, data...)

	if len(clientIDs) == 0 {
		for _, client := range server.clients {
			_, err := client.Write(buffer)
			if err != nil {
				return err
			}
		}
	} else {
		for _, clientID := range clientIDs {
			if client, ok := server.clients[clientID]; ok {
				_, err := client.Write(buffer)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("client does not exist")
			}
		}
	}

	return nil
}

// Serving returns a boolean value representing whether or not the server is
// serving
func (server *Server) Serving() bool {
	return server.serving
}

// GetAddr returns the address string
func (server *Server) GetAddr() (string, uint16, error) {
	if !server.serving {
		return "", 0, fmt.Errorf("server is not serving")
	}

	return parseAddr(server.sock.Addr().String())
}

// GetClientAddr returns the address of a client
func (server *Server) GetClientAddr(clientID uint) (string, uint16, error) {
	if !server.serving {
		return "", 0, fmt.Errorf("server is not serving")
	}

	if client, ok := server.clients[clientID]; ok {
		return parseAddr(client.RemoteAddr().String())
	}
	return "", 0, fmt.Errorf("client does not exist")
}

// RemoveClient disconnects a client from the server
func (server *Server) RemoveClient(clientID uint) error {
	if !server.serving {
		return fmt.Errorf("server is not serving")
	}

	if client, ok := server.clients[clientID]; ok {
		client.Close()
		delete(server.clients, clientID)
		return nil
	}
	return fmt.Errorf("client does not exist")
}

// Handle client connections
func (server *Server) serve() {
	defer server.wg.Done()

	for ; server.serving; {
		conn, err := server.sock.Accept()
		if err != nil {
			if server.serving {
				fmt.Println("ACCEPT ERROR:", err)
			}
			break
		}

		clientID := server.newClientID()
		server.clients[clientID] = conn
		server.wg.Add(1)
		go server.serveClient(clientID)
	}
}

// Serve clients
func (server *Server) serveClient(clientID uint) {
	defer server.wg.Done()

	if server.eventBlocking {
		if server.onConnect != nil {
			server.onConnect(clientID)
		}
		if server.onDisconnect != nil {
			defer server.onDisconnect(clientID)
		}
	} else {
		if server.onConnect != nil {
			go server.onConnect(clientID)
		}
		if server.onDisconnect != nil {
			defer func() {
				go server.onDisconnect(clientID)
			}()
		}
	}

	client := server.clients[clientID]

	sizebuffer := make([]byte, lenSize)
	for ; server.serving; {
		_, err := client.Read(sizebuffer)
		if err != nil {
			break
		}

		msgSize := asciiToDec(sizebuffer)
		buffer := make([]byte, msgSize)
		_, err = client.Read(buffer)
		if err != nil {
			break
		}

		if server.onRecv != nil {
			if server.eventBlocking {
				server.onRecv(clientID, buffer)
			} else {
				go server.onRecv(clientID, buffer)
			}
		}
	}
}

// Get a new client ID
func (server *Server) newClientID() uint {
	server.nextClientID++
	return server.nextClientID - 1
}
