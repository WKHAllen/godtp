# Go Data Transfer Protocol

GoDTP is a cross platform networking library written in Go. It is based on [dtplib](https://github.com/WKHAllen/dtplib) and [cdtp](https://github.com/WKHAllen/cdtp).

## Installation

First, install the package:

```console
$ go get -u github.com/WKHAllen/godtp
```

Then import it:

```go
import "github.com/WKHAllen/godtp"
```

## Example

See [the test file](godtp_test.go) for a basic example.

## Encoding and Decoding

GoDTP provides functions for encoding and decoding objects. This is necessary, as only bytes can be sent over a network.

### func Decode

```go
func Decode(bytestring []byte, object interface{}) error
```

Decode an object that came through a socket.

### func Encode

```go
func Encode(object interface{}) ([]byte, error)
```

Encode an object so it can be sent through a socket.

## Client (using callbacks)

Below are the defined client types and functions:

### type Client

```go
type Client struct {
    // contains filtered or unexported fields
}
```

Client defines the socket client type.

### func NewClient

```go
func NewClient(onRecv onRecvFuncClient,
    onDisconnected onDisconnectedFuncClient,
    blocking bool, eventBlocking bool) *Client
```

NewClient creates a new socket client object.

### func NewClientDefault

```go
func NewClientDefault(onRecv onRecvFuncClient,
    onDisconnected onDisconnectedFuncClient) *Client
```

NewClientDefault creates a new socket client object with blocking and eventBlocking set to false.

### func (*Client) Connect

```go
func (client *Client) Connect(host string, port uint16) error
```

Connect to a server.

### func (*Client) ConnectDefaultHost

```go
func (client *Client) ConnectDefaultHost(port uint16) error
```

ConnectDefaultHost connects to a server at the default host address.

### func (*Client) Connected

```go
func (client *Client) Connected() bool
```

Connected returns a boolean value representing whether or not the client is connected to a server.

### func (*Client) Disconnect

```go
func (client *Client) Disconnect() error
```

Disconnect from the server.

### func (*Client) GetAddr

```go
func (client *Client) GetAddr() (string, uint16, error)
```

GetAddr returns the address string.

### func (*Client) GetServerAddr

```go
func (client *Client) GetServerAddr() (string, uint16, error)
```

GetServerAddr returns the address of the server.

### func (*Client) Send

```go
func (client *Client) Send(data []byte) error
```

Send data to the server.

## Server (using callbacks)

Below are the defined server types and functions:

### type Server

```go
type Server struct {
    // contains filtered or unexported fields
}
```

Server defines the socket server type.

### func NewServer

```go
func NewServer(onRecv onRecvFuncServer,
    onConnect onConnectFuncServer,
    onDisconnect onDisconnectFuncServer,
    blocking bool, eventBlocking bool) *Server
```

NewServer creates a new socket server object.

### func NewServerDefault

```go
func NewServerDefault(onRecv onRecvFuncServer,
    onConnect onConnectFuncServer,
    onDisconnect onDisconnectFuncServer) *Server
```

NewServerDefault creates a new socket server object with blocking and eventBlocking set to false.

### func (*Server) GetAddr

```go
func (server *Server) GetAddr() (string, uint16, error)
```

GetAddr returns the address string.

### func (*Server) GetClientAddr

```go
func (server *Server) GetClientAddr(clientID uint) (string, uint16, error)
```

GetClientAddr returns the address of a client.

### func (*Server) RemoveClient

```go
func (server *Server) RemoveClient(clientID uint) error
```

RemoveClient disconnects a client from the server.

### func (*Server) Send

```go
func (server *Server) Send(data []byte, clientIDs ...uint) error
```

Send data to clients.

### func (*Server) Serving

```go
func (server *Server) Serving() bool
```

Serving returns a boolean value representing whether or not the server is serving.

### func (*Server) Start

```go
func (server *Server) Start(host string, port uint16) error
```

Start the server.

### func (*Server) StartDefault

```go
func (server *Server) StartDefault() error
```

StartDefault starts the server on 0.0.0.0:0.

### func (*Server) StartDefaultHost

```go
func (server *Server) StartDefaultHost(port uint16) error
```

StartDefaultHost starts the server at the default host address.

### func (*Server) StartDefaultPort

```go
func (server *Server) StartDefaultPort(host string) error
```

StartDefaultPort starts the server on the default port.

### func (*Server) Stop

```go
func (server *Server) Stop() error
```

Stop the server.

## Client (using channels)

Below are the defined client channel types and functions:

### type ClientChan

```go
type ClientChan struct {
    // contains filtered or unexported fields
}
```

ClientChan defines the socket client chan type.

### func NewClientChan

```go
func NewClientChan() *ClientChan
```

NewClientChan creates a new socket client object, using channels rather than callbacks.

### func (*ClientChan) Connect

```go
func (client *ClientChan) Connect(host string, port uint16) (chan []byte, chan []byte, error)
```

Connect to a server.

### func (*ClientChan) ConnectDefaultHost

```go
func (client *ClientChan) ConnectDefaultHost(port uint16) (chan []byte, chan []byte, error)
```

ConnectDefaultHost connects to a server at the default host address.

### func (*ClientChan) Connected

```go
func (client *ClientChan) Connected() bool
```

Connected returns a boolean value representing whether or not the client is connected to a server.

### func (*ClientChan) Disconnect

```go
func (client *ClientChan) Disconnect() error
```

Disconnect from the server.

### func (*ClientChan) GetAddr

```go
func (client *ClientChan) GetAddr() (string, uint16, error)
```

GetAddr returns the address string.

### func (*ClientChan) GetServerAddr

```go
func (client *ClientChan) GetServerAddr() (string, uint16, error)
```

GetServerAddr returns the address of the server.

## Server (using channels)

Below are the defined server channel types and functions:

### type ServerChan

```go
type ServerChan struct {
    // contains filtered or unexported fields
}
```

ServerChan defines the socket server chan type.

### func NewServerChan

```go
func NewServerChan() *ServerChan
```

NewServerChan creates a new socket server object, using channels rather than callbacks.

### func (*ServerChan) GetAddr

```go
func (server *ServerChan) GetAddr() (string, uint16, error)
```

GetAddr returns the address string.

### func (*ServerChan) GetClientAddr

```go
func (server *ServerChan) GetClientAddr(clientID uint) (string, uint16, error)
```

GetClientAddr returns the address of a client.

### func (*ServerChan) Serving

```go
func (server *ServerChan) Serving() bool
```

Serving returns a boolean value representing whether or not the server is serving.

### func (*ServerChan) Start

```go
func (server *ServerChan) Start(host string, port uint16) (chan chan []byte, error)
```

Start the server.

### func (*ServerChan) StartDefault

```go
func (server *ServerChan) StartDefault() (chan chan []byte, error)
```

StartDefault starts the server on 0.0.0.0:0.

### func (*ServerChan) StartDefaultHost

```go
func (server *ServerChan) StartDefaultHost(port uint16) (chan chan []byte, error)
```

StartDefaultHost starts the server at the default host address.

### func (*ServerChan) StartDefaultPort

```go
func (server *ServerChan) StartDefaultPort(host string) (chan chan []byte, error)
```

StartDefaultPort starts the server on the default port.

### func (*ServerChan) Stop

```go
func (server *ServerChan) Stop() error
```

Stop the server.
