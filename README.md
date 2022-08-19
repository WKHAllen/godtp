# Go Data Transfer Protocol

Cross-platform networking interfaces for Go.

## Data Transfer Protocol

The Data Transfer Protocol (DTP) is a larger project to make ergonomic network programming available in any language.
See the full project [here](https://wkhallen.com/dtp/).

## Installation

Install the package:

```console
$ go get -u github.com/wkhallen/godtp
```

## Creating a server

A server can be built using the `Server` implementation:

```go
package example

import (
	"fmt"
	"github.com/wkhallen/godtp"
)

func main() {
	// Create a server that receives strings and returns the length of each string
	server, serverEvent := godtp.NewServer[int, string]()
	err := server.Start("127.0.0.1", 29275)
	if err != nil {
		// Handle server start error
	}

	// Iterate over events
	for event := range serverEvent {
		switch event.EventType {
		case godtp.ServerConnect:
			fmt.Printf("Client with ID %d connected\n", event.ClientID)
		case godtp.ServerDisconnect:
			fmt.Printf("Client with ID %d disconnected\n", event.ClientID)
		case godtp.ServerReceive:
			// Send back the length of the string
			err := server.Send(len(event.Data))
			if err != nil {
				// Handle send error
			}
		}
	}
}
```

## Creating a client

A client can be built using the `Client` implementation:

```go
package example

import (
	"fmt"
	"github.com/wkhallen/godtp"
)

func main() {
	// Create a client that send a message to the server and receives the length of the message
	client, clientEvent := godtp.NewClient[string, int]()
	err := client.Connect("127.0.0.1", 29275)
	if err != nil {
		// Handle client connect error
	}

	// Send a message to the server
	message := "Hello, server!"
	err = client.Send(message)
	if err != nil {
		// Handle send error
	}

	// Receive the response
	event := <-clientEvent
	switch event.EventType {
	case godtp.ClientReceive:
		// Validate the response
		fmt.Printf("Received response from server: %d", event.Data)
		if event.Data != len(message) {
			fmt.Errorf("invalid response: expected %d, received %d", len(message), event.Data)
		}
	default:
		// Unexpected response
		fmt.Errorf("expected to receive a response from the server, instead got %#v\n", event)
	}
}
```

## Security

Information security comes included. Every message sent over a network interface is encrypted with AES-256. Key
exchanges are performed using a 2048-bit RSA key-pair.
