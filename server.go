package godtp

import (
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"sync"
)

// ServerEventType defines the type of server event
type ServerEventType uint

// Server event type values
const (
	ServerReceive ServerEventType = iota
	ServerConnect
	ServerDisconnect
)

// ServerEvent defines an event emitted from the server
type ServerEvent[T any] struct {
	EventType ServerEventType
	ClientID  uint
	Data      T
}

// Server defines the socket server type
type Server[S any, R any] struct {
	serving      bool
	sock         net.Listener
	clients      map[uint]net.Conn
	keys         map[uint][]byte
	eventChannel chan ServerEvent[R]
	wg           sync.WaitGroup
	nextClientID uint
}

// NewServer creates a new socket server
func NewServer[S any, R any]() (*Server[S, R], <-chan ServerEvent[R]) {
	eventChannel := make(chan ServerEvent[R], channelBufferSize)

	return &Server[S, R]{
		serving:      false,
		clients:      make(map[uint]net.Conn),
		keys:         make(map[uint][]byte),
		eventChannel: eventChannel,
		nextClientID: 0,
	}, eventChannel
}

// Start the server
func (server *Server[S, R]) Start(host string, port uint16) error {
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
	server.wg.Add(1)
	go server.serve()

	return nil
}

// Stop the server
func (server *Server[S, R]) Stop() error {
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

	server.wg.Wait()
	close(server.eventChannel)

	return nil
}

// Send data to clients
func (server *Server[S, R]) Send(data S, clientIDs ...uint) error {
	if !server.serving {
		return fmt.Errorf("server is not serving")
	}

	dataBytes, err := encodeObject(data)
	if err != nil {
		return err
	}

	if len(clientIDs) == 0 {
		for clientID := range server.clients {
			clientIDs = append(clientIDs, clientID)
		}
	}

	for _, clientID := range clientIDs {
		if client, ok := server.clients[clientID]; ok {
			encryptedData, err := aesEncrypt(server.keys[clientID], dataBytes)
			if err != nil {
				return err
			}

			size := encodeMessageSize(uint64(len(encryptedData)))
			buffer := append(size, encryptedData...)

			_, err = client.Write(buffer)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("client does not exist")
		}
	}

	return nil
}

// Serving returns a boolean value representing whether the server is serving
func (server *Server[S, R]) Serving() bool {
	return server.serving
}

// GetAddr returns the server's address
func (server *Server[S, R]) GetAddr() (string, uint16, error) {
	if !server.serving {
		return "", 0, fmt.Errorf("server is not serving")
	}

	return parseAddr(server.sock.Addr().String())
}

// GetClientAddr returns a client's address
func (server *Server[S, R]) GetClientAddr(clientID uint) (string, uint16, error) {
	if !server.serving {
		return "", 0, fmt.Errorf("server is not serving")
	}

	if client, ok := server.clients[clientID]; ok {
		return parseAddr(client.RemoteAddr().String())
	}
	return "", 0, fmt.Errorf("client does not exist")
}

// RemoveClient disconnects a client from the server
func (server *Server[S, R]) RemoveClient(clientID uint) error {
	if !server.serving {
		return fmt.Errorf("server is not serving")
	}

	if client, ok := server.clients[clientID]; ok {
		err := client.Close()
		if err != nil {
			return err
		}

		delete(server.clients, clientID)
		delete(server.keys, clientID)

		return nil
	}
	return fmt.Errorf("client does not exist")
}

// Handle client connections
func (server *Server[S, R]) serve() {
	defer func() {
		server.wg.Done()
	}()

	for server.serving {
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
func (server *Server[S, R]) serveClient(clientID uint) {
	defer func() {
		server.wg.Done()
	}()

	server.eventChannel <- ServerEvent[R]{
		EventType: ServerConnect,
		ClientID:  clientID,
	}
	defer func() {
		server.eventChannel <- ServerEvent[R]{
			EventType: ServerDisconnect,
			ClientID:  clientID,
		}
	}()

	client := server.clients[clientID]

	sizeBuffer := make([]byte, lenSize)

	for server.serving {
		_, err := client.Read(sizeBuffer)
		if err != nil {
			break
		}

		msgSize := decodeMessageSize(sizeBuffer)
		buffer := make([]byte, msgSize)
		_, err = client.Read(buffer)
		if err != nil {
			break
		}

		dataBytes, err := aesDecrypt(server.keys[clientID], buffer)
		if err != nil {
			break
		}

		data, err := decodeObject[R](dataBytes)
		if err != nil {
			break
		}

		server.eventChannel <- ServerEvent[R]{
			EventType: ServerReceive,
			ClientID:  clientID,
			Data:      data,
		}
	}
}

// Get a new client ID
func (server *Server[S, R]) newClientID() uint {
	server.nextClientID++
	return server.nextClientID - 1
}

// Exchange crypto keys with a client
func (server *Server[S, R]) exchangeKeys(clientID uint, client net.Conn) error {
	privateKey, err := newRSAKeys()
	if err != nil {
		return err
	}

	pub := privateKey.PublicKey
	enc := gob.NewEncoder(client)
	err = enc.Encode(&pub)
	if err != nil {
		return err
	}

	sizeBuffer := make([]byte, lenSize)
	_, err = client.Read(sizeBuffer)
	if err != nil {
		return err
	}

	msgSize := decodeMessageSize(sizeBuffer)
	buffer := make([]byte, msgSize)
	_, err = client.Read(buffer)
	if err != nil {
		return err
	}

	key, err := rsaDecrypt(privateKey, buffer)
	if err != nil {
		return err
	}

	server.keys[clientID] = key

	return nil
}
