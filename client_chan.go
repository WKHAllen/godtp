package godtp

// ClientChan defines the socket client chan type
type ClientChan struct {
	client *Client
	sendChan chan []byte
	recvChan chan []byte
}

// NewClientChan creates a new socket client object, using channels rather
// than callbacks
func NewClientChan() *ClientChan {
	client := &ClientChan{
		client: NewClient(nil, nil, false, false),
		sendChan: make(chan []byte),
		recvChan: make(chan []byte),
	}
	client.client.onRecv = client.onRecvCallback
	client.client.onDisconnected = client.onDisconnectedCallback
	return client
}

// Connect to a server
func (client *ClientChan) Connect(host string, port uint16) (chan []byte, chan []byte, error) {
	err := client.client.Connect(host, port)
	if err != nil {
		return nil, nil, err
	}

	go client.handle()

	return client.sendChan, client.recvChan, nil
}

// ConnectDefaultHost connects to a server at the default host address
func (client *ClientChan) ConnectDefaultHost(port uint16) (chan []byte, chan []byte, error) {
	return client.Connect("0.0.0.0", port)
}

// Disconnect from the server
func (client *ClientChan) Disconnect() error {
	return client.client.Disconnect()
}

// Connected returns a boolean value representing whether or not the client is
// connected to a server
func (client *ClientChan) Connected() bool {
	return client.client.Connected()
}

// GetAddr returns the address string
func (client *ClientChan) GetAddr() (string, uint16, error) {
	return client.client.GetAddr()
}

// GetServerAddr returns the address of the server
func (client *ClientChan) GetServerAddr() (string, uint16, error) {
	return client.client.GetServerAddr()
}

// Handle message sending and channel closing
func (client *ClientChan) handle() {
	for msg := range client.sendChan {
		client.client.Send(msg)
	}

	client.onDisconnectedCallback()
}

// Handle message receiving
func (client *ClientChan) onRecvCallback(msg []byte) {
	client.recvChan <- msg
}

// Handle server disconnecting
func (client *ClientChan) onDisconnectedCallback() {
	if _, ok := <-client.sendChan; ok {
		close(client.sendChan)
	}
	if _, ok := <-client.recvChan; ok {
		close(client.recvChan)
	}
}
