package godtp

import (
	"fmt"
	"net"
	"strconv"
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

// Connect to a server
func (client *Client) Connect(host string, port uint16) error {
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
	if client.blocking {
		client.handle()
	} else {
		client.wg.Add(1)
		go client.handle()
	}

	return nil
}

// ConnectDefaultHost connects to a server at the default host address
func (client *Client) ConnectDefaultHost(port uint16) error {
	return client.Connect("0.0.0.0", port)
}

// Disconnect from the server
func (client *Client) Disconnect() error {
	if !client.connected {
		return fmt.Errorf("client is not connected to a server")
	}

	client.connected = false

	err := client.sock.Close()
	if err != nil {
		return err
	}

	if !client.blocking {
		client.wg.Wait()
	}

	return nil
}

// Send data to the server
func (client *Client) Send(data []byte) error {
	if !client.connected {
		return fmt.Errorf("client is not connected to a server")
	}

	size := decToASCII(uint64(len(data)))
	buffer := append(size, data...)

	_, err := client.sock.Write(buffer)
	if err != nil {
		return err
	}

	return nil
}

// Connected returns a boolean value representing whether or not the client is
// connected to a server
func (client *Client) Connected() bool {
	return client.connected
}

// GetAddr returns the address string
func (client *Client) GetAddr() (string, uint16, error) {
	if !client.connected {
		return "", 0, fmt.Errorf("client is not connected to a server")
	}

	return parseAddr(client.sock.LocalAddr().String())
}

// GetServerAddr returns the address of the server
func (client *Client) GetServerAddr() (string, uint16, error) {
	if !client.connected {
		return "", 0, fmt.Errorf("client is not connected to a server")
	}

	return parseAddr(client.sock.RemoteAddr().String())
}

// Handle client events
func (client *Client) handle() {
	defer client.wg.Done()

	sizebuffer := make([]byte, lenSize)
	for ; client.connected; {
		_, err := client.sock.Read(sizebuffer)
		if err != nil {
			break
		}

		msgSize := asciiToDec(sizebuffer)
		buffer := make([]byte, msgSize)
		_, err = client.sock.Read(buffer)
		if err != nil {
			break
		}

		if client.eventBlocking {
			client.onRecv(buffer)
		} else {
			go client.onRecv(buffer)
		}
	}

	if client.connected {
		client.connected = false
		client.sock.Close()
		if client.eventBlocking {
			client.onDisconnected()
		} else {
			go client.onDisconnected()
		}
	}
}
