package godtp

import (
	"crypto/rsa"
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"sync"
)

// ClientEventType defines the type of client event
type ClientEventType uint

// Client event type values
const (
	ClientReceive ClientEventType = iota
	ClientDisconnected
)

// ClientEvent defines an event emitted from the client
type ClientEvent[T any] struct {
	EventType ClientEventType
	Data      T
}

// Client defines the socket client type
type Client[S any, R any] struct {
	connected    bool
	sock         net.Conn
	key          []byte
	eventChannel chan ClientEvent[R]
	wg           sync.WaitGroup
}

// NewClient creates a new socket client
func NewClient[S any, R any]() (*Client[S, R], <-chan ClientEvent[R]) {
	eventChannel := make(chan ClientEvent[R], channelBufferSize)

	return &Client[S, R]{
		connected:    false,
		eventChannel: eventChannel,
	}, eventChannel
}

// Connect to a server
func (client *Client[S, R]) Connect(host string, port uint16) error {
	if client.connected {
		return fmt.Errorf("client is already connected to a server")
	}

	address := host + ":" + strconv.Itoa(int(port))
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	client.sock = conn
	client.connected = true

	err = client.exchangeKeys()
	if err != nil {
		return err
	}

	client.wg.Add(1)
	go client.handle()

	return nil
}

// Disconnect from the server
func (client *Client[S, R]) Disconnect() error {
	if !client.connected {
		return fmt.Errorf("client is not connected to a server")
	}

	client.connected = false

	err := client.sock.Close()
	if err != nil {
		return err
	}

	client.wg.Wait()
	close(client.eventChannel)

	return nil
}

// Send data to the server
func (client *Client[S, R]) Send(data S) error {
	if !client.connected {
		return fmt.Errorf("client is not connected to a server")
	}

	dataBytes, err := encodeObject(data)
	if err != nil {
		return err
	}

	encryptedData, err := aesEncrypt(client.key, dataBytes)
	if err != nil {
		return err
	}

	size := encodeMessageSize(uint64(len(encryptedData)))
	buffer := append(size, encryptedData...)

	_, err = client.sock.Write(buffer)
	if err != nil {
		return err
	}

	return nil
}

// Connected returns a boolean value representing whether the client is connected to a server
func (client *Client[S, R]) Connected() bool {
	return client.connected
}

// GetAddr returns the client's address
func (client *Client[S, R]) GetAddr() (string, uint16, error) {
	if !client.connected {
		return "", 0, fmt.Errorf("client is not connected to a server")
	}

	return parseAddr(client.sock.LocalAddr().String())
}

// GetServerAddr returns the server's address
func (client *Client[S, R]) GetServerAddr() (string, uint16, error) {
	if !client.connected {
		return "", 0, fmt.Errorf("client is not connected to a server")
	}

	return parseAddr(client.sock.RemoteAddr().String())
}

// Handle client events
func (client *Client[S, R]) handle() {
	defer client.wg.Done()

	sizeBuffer := make([]byte, lenSize)

	for client.connected {
		_, err := client.sock.Read(sizeBuffer)
		if err != nil {
			break
		}

		msgSize := decodeMessageSize(sizeBuffer)
		buffer := make([]byte, msgSize)
		_, err = client.sock.Read(buffer)
		if err != nil {
			break
		}

		dataBytes, err := aesDecrypt(client.key, buffer)
		if err != nil {
			break
		}

		data, err := decodeObject[R](dataBytes)
		if err != nil {
			break
		}

		client.eventChannel <- ClientEvent[R]{
			EventType: ClientReceive,
			Data:      data,
		}
	}

	if client.connected {
		client.connected = false
		err := client.sock.Close()
		if err != nil {
			// Do nothing for these errors
		}

		client.eventChannel <- ClientEvent[R]{
			EventType: ClientDisconnected,
		}
	}
}

// Exchange crypto keys with the server
func (client *Client[S, R]) exchangeKeys() error {
	pub := rsa.PublicKey{}
	dec := gob.NewDecoder(client.sock)
	err := dec.Decode(&pub)
	if err != nil {
		return err
	}

	key, err := newAESKey()
	if err != nil {
		return err
	}

	encryptedKey, err := rsaEncrypt(pub, key)

	size := encodeMessageSize(uint64(len(encryptedKey)))
	buffer := append(size, encryptedKey...)

	_, err = client.sock.Write(buffer)
	if err != nil {
		return err
	}

	client.key = key

	return nil
}
