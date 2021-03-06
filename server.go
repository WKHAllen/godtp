package godtp

import (
	"crypto/rsa"
	"encoding/gob"
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
	keys map[uint][]byte
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
		keys: make(map[uint][]byte),
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

	if len(clientIDs) == 0 {
		for clientID, client := range server.clients {
			encryptedData, err := encrypt(server.keys[clientID], data)
			if err != nil {
				return err
			}

			size := decToASCII(uint64(len(encryptedData)))
			buffer := append(size, encryptedData...)

			_, err = client.Write(buffer)
			if err != nil {
				return err
			}
		}
	} else {
		for _, clientID := range clientIDs {
			if client, ok := server.clients[clientID]; ok {
				encryptedData, err := encrypt(server.keys[clientID], data)
				if err != nil {
					return err
				}

				size := decToASCII(uint64(len(encryptedData)))
				buffer := append(size, encryptedData...)

				_, err = client.Write(buffer)
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
		delete(server.keys, clientID)
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
			if !server.serving {
				break
			} else {
				continue
			}
		}

		clientID := server.newClientID()
		err = server.exchangeKeys(clientID, conn)
		if err != nil {
			if !server.serving {
				break
			} else {
				continue
			}
		}

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

		data, err := decrypt(server.keys[clientID], buffer)
		if err != nil {
			break
		}

		if server.onRecv != nil {
			if server.eventBlocking {
				server.onRecv(clientID, data)
			} else {
				go server.onRecv(clientID, data)
			}
		}
	}
}

// Get a new client ID
func (server *Server) newClientID() uint {
	server.nextClientID++
	return server.nextClientID - 1
}

// Exchange keys with a client
func (server *Server) exchangeKeys(clientID uint, client net.Conn) error {
	pub := rsa.PublicKey{}
	dec := gob.NewDecoder(client)
	dec.Decode(&pub)

	key, err := newKey()
	if err != nil {
		return err
	}

	encryptedKey, err := asymmetricEncrypt(pub, key)

	size := decToASCII(uint64(len(encryptedKey)))
	buffer := append(size, encryptedKey...)

	_, err = client.Write(buffer)
	if err != nil {
		return err
	}

	server.keys[clientID] = key

	return nil
}
