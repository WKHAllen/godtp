package godtp

// ServerChan defines the socket server chan type
type ServerChan struct {
	server *Server
	sendChans map[uint]chan []byte
	recvChans map[uint]chan []byte
	connectChan chan chan []byte
}

// NewServerChan creates a new socket server object, using channels rather
// than callbacks
func NewServerChan() *ServerChan {
	server := &ServerChan{
		server: NewServer(nil, nil, nil, false, false),
		sendChans: make(map[uint]chan []byte),
		recvChans: make(map[uint]chan []byte),
		connectChan: make(chan chan []byte, 2),
	}
	server.server.onRecv = server.onRecvCallback
	server.server.onConnect = server.onConnectCallback
	server.server.onDisconnect = server.onDisconnectCallback
	return server
}

// Start the server
func (server *ServerChan) Start(host string, port uint16) (chan chan []byte, error) {
	err := server.server.Start(host, port)
	if err != nil {
		return nil, err
	}
	return server.connectChan, nil
}

// StartDefaultHost starts the server at the default host address
func (server *ServerChan) StartDefaultHost(port uint16) (chan chan []byte, error) {
	return server.Start("0.0.0.0", port)
}

// StartDefaultPort starts the server on the default port
func (server *ServerChan) StartDefaultPort(host string) (chan chan []byte, error) {
	return server.Start(host, 0)
}

// StartDefault starts the server on 0.0.0.0:0
func (server *ServerChan) StartDefault() (chan chan []byte, error) {
	return server.Start("0.0.0.0", 0)
}

// Stop the server
func (server *ServerChan) Stop() error {
	return server.server.Stop()
}

// Serving returns a boolean value representing whether or not the server is
// serving
func (server *ServerChan) Serving() bool {
	return server.server.Serving()
}

// GetAddr returns the address string
func (server *ServerChan) GetAddr() (string, uint16, error) {
	return server.server.GetAddr()
}

// GetClientAddr returns the address of a client
func (server *ServerChan) GetClientAddr(clientID uint) (string, uint16, error) {
	return server.server.GetClientAddr(clientID)
}

// Handle message sending, channel closing, and client disconnecting
func (server *ServerChan) serveClient(clientID uint) {
	for msg := range server.sendChans[clientID] {
		server.server.Send(msg, clientID)
	}

	server.onDisconnectCallback(clientID)
	server.server.RemoveClient(clientID)
}

// Handle message receiving
func (server *ServerChan) onRecvCallback(clientID uint, msg []byte) {
	server.recvChans[clientID] <- msg
}

// Handle client connecting
func (server *ServerChan) onConnectCallback(clientID uint) {
	sendChan := make(chan []byte)
	recvChan := make(chan []byte)
	server.sendChans[clientID] = sendChan
	server.recvChans[clientID] = recvChan

	go server.serveClient(clientID)

	server.connectChan <- sendChan
	server.connectChan <- recvChan
}

// Handle client disconnecting
func (server *ServerChan) onDisconnectCallback(clientID uint) {
	if _, ok := <-server.sendChans[clientID]; ok {
		close(server.sendChans[clientID])
	}
	if _, ok := <-server.recvChans[clientID]; ok {
		close(server.recvChans[clientID])
	}
	delete(server.sendChans, clientID)
	delete(server.recvChans, clientID)
}
