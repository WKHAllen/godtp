package godtp

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Amount of time to wait between network operations
const waitTime = 100 * time.Millisecond

// Assert a condition
func assert(value bool, t *testing.T, err string) {
	if !value {
		t.Errorf(err)
		panic(err)
	}
}

// Assert an equality
func assertEq[T any](left T, right T, t *testing.T) {
	if !reflect.DeepEqual(left, right) {
		t.Errorf("Assertion error: %#v != %#v", left, right)
		panic("Assertion error")
	}
}

// Assert an inequality
func assertNe[T any](left T, right T, t *testing.T) {
	if reflect.DeepEqual(left, right) {
		t.Errorf("Assertion error: %#v == %#v", left, right)
		panic("Assertion error")
	}
}

// Assert no error occurred
func assertNoErr(err error, t *testing.T) {
	if err != nil {
		t.Errorf(err.Error())
		panic(err.Error())
	}
}

// A representation of a person
type person struct {
	Name        string
	Age         int
	WritesInGo  bool
	PrefersRust bool
}

// Custom type
type custom struct {
	A int
	B string
	C []string
}

// Test encoding message sizes
func TestEncodeMessageSize(t *testing.T) {
	assertEq(encodeMessageSize(0), []byte{0, 0, 0, 0, 0}, t)
	assertEq(encodeMessageSize(1), []byte{0, 0, 0, 0, 1}, t)
	assertEq(encodeMessageSize(255), []byte{0, 0, 0, 0, 255}, t)
	assertEq(encodeMessageSize(256), []byte{0, 0, 0, 1, 0}, t)
	assertEq(encodeMessageSize(257), []byte{0, 0, 0, 1, 1}, t)
	assertEq(encodeMessageSize(4311810305), []byte{1, 1, 1, 1, 1}, t)
	assertEq(encodeMessageSize(4328719365), []byte{1, 2, 3, 4, 5}, t)
	assertEq(encodeMessageSize(47362409218), []byte{11, 7, 5, 3, 2}, t)
	assertEq(encodeMessageSize(1099511627775), []byte{255, 255, 255, 255, 255}, t)
}

// Test decoding message sizes
func TestDecodeMessageSize(t *testing.T) {
	assertEq(decodeMessageSize([]byte{0, 0, 0, 0, 0}), 0, t)
	assertEq(decodeMessageSize([]byte{0, 0, 0, 0, 1}), 1, t)
	assertEq(decodeMessageSize([]byte{0, 0, 0, 0, 255}), 255, t)
	assertEq(decodeMessageSize([]byte{0, 0, 0, 1, 0}), 256, t)
	assertEq(decodeMessageSize([]byte{0, 0, 0, 1, 1}), 257, t)
	assertEq(decodeMessageSize([]byte{1, 1, 1, 1, 1}), 4311810305, t)
	assertEq(decodeMessageSize([]byte{1, 2, 3, 4, 5}), 4328719365, t)
	assertEq(decodeMessageSize([]byte{11, 7, 5, 3, 2}), 47362409218, t)
	assertEq(decodeMessageSize([]byte{255, 255, 255, 255, 255}), 1099511627775, t)
}

// Test object encoding and decoding
func TestObjectEncodeDecode(t *testing.T) {
	// Encode/decode an integer
	intValue := 29275
	encodedInt, err := encodeObject(intValue)
	assertNoErr(err, t)
	decodedInt, err := decodeObject[int](encodedInt)
	assertNoErr(err, t)
	assertEq(decodedInt, intValue, t)

	// Encode/decode a string
	stringValue := "Hello, encoder!"
	encodedString, err := encodeObject(stringValue)
	assertNoErr(err, t)
	decodedString, err := decodeObject[string](encodedString)
	assertNoErr(err, t)
	assertEq(decodedString, stringValue, t)

	// Encode/decode a slice
	sliceValue := []int{2, 3, 5, 7, 11}
	encodedSlice, err := encodeObject(sliceValue)
	assertNoErr(err, t)
	decodedSlice, err := decodeObject[[]int](encodedSlice)
	assertNoErr(err, t)
	assertEq(decodedSlice, sliceValue, t)

	// Encode/decode a map
	mapValue := map[int]int{0: 1, 1: 1, 2: 2, 3: 6, 4: 24, 5: 120}
	encodedMap, err := encodeObject(mapValue)
	assertNoErr(err, t)
	decodedMap, err := decodeObject[map[int]int](encodedMap)
	assertNoErr(err, t)
	assertEq(decodedMap, mapValue, t)

	// Encode/decode a struct
	structValue := person{
		Name:        "Will",
		Age:         24,
		WritesInGo:  true,
		PrefersRust: true,
	}
	encodedStruct, err := encodeObject(structValue)
	assertNoErr(err, t)
	decodedStruct, err := decodeObject[person](encodedStruct)
	assertNoErr(err, t)
	assertEq(decodedStruct, structValue, t)
}

// Test crypto functions
func TestCrypto(t *testing.T) {
	// Test RSA encryption
	rsaMessage := "Hello, RSA!"
	privateKey, err := newRSAKeys()
	assertNoErr(err, t)
	publicKey := privateKey.PublicKey
	rsaEncrypted, err := rsaEncrypt(publicKey, []byte(rsaMessage))
	assertNoErr(err, t)
	rsaDecrypted, err := rsaDecrypt(privateKey, rsaEncrypted)
	assertNoErr(err, t)
	rsaDecryptedMessage := string(rsaDecrypted[:])
	assertEq(rsaDecryptedMessage, rsaMessage, t)
	assertNe(rsaEncrypted, []byte(rsaMessage), t)

	// Test AES encryption
	aesMessage := "Hello, AES!"
	key, err := newAESKey()
	assertNoErr(err, t)
	aesEncrypted, err := aesEncrypt(key, []byte(aesMessage))
	assertNoErr(err, t)
	aesDecrypted, err := aesDecrypt(key, aesEncrypted)
	assertNoErr(err, t)
	aesDecryptedMessage := string(aesDecrypted[:])
	assertEq(aesDecryptedMessage, aesMessage, t)
	assertNe(aesEncrypted, []byte(aesMessage), t)

	// Test encrypting an AES key with RSA
	encryptedKey, err := rsaEncrypt(publicKey, key)
	assertNoErr(err, t)
	decryptedKey, err := rsaDecrypt(privateKey, encryptedKey)
	assertNoErr(err, t)
	assertEq(decryptedKey, key, t)
	assertNe(encryptedKey, key, t)
}

// Test server creation and serving
func TestServerServe(t *testing.T) {
	// Create server
	server, _ := NewServer[any, any]()
	assert(!server.Serving(), t, "Server should not be serving")

	// Start server
	err := server.Start("0.0.0.0", 0)
	assertNoErr(err, t)
	assert(server.Serving(), t, "Server should be serving")
	time.Sleep(waitTime)

	// Check server address info
	host, port, err := server.GetAddr()
	assertNoErr(err, t)
	assert(server.sock.Addr().String() == host+":"+strconv.Itoa(int(port)), t, "Address strings don't match")
	fmt.Printf("Server address: %s:%d\n", host, port)

	// Stop server
	err = server.Stop()
	assertNoErr(err, t)
	assert(!server.Serving(), t, "Server should not be serving")
	time.Sleep(waitTime)
}

// Test getting server and client addresses
func TestAddresses(t *testing.T) {
	// Create server
	server, serverEvent := NewServer[any, any]()
	assert(!server.Serving(), t, "Server should not be serving")

	// Start server
	err := server.Start("127.0.0.1", 0)
	assertNoErr(err, t)
	assert(server.Serving(), t, "Server should be serving")
	time.Sleep(waitTime)

	// Check server address info
	host, port, err := server.GetAddr()
	assertNoErr(err, t)
	assert(server.sock.Addr().String() == host+":"+strconv.Itoa(int(port)), t, "Address strings don't match")
	fmt.Printf("Server address: %s:%d\n", host, port)

	// Create client
	client, _ := NewClient[any, any]()
	assert(!client.Connected(), t, "Client should not be connected")

	// Connect to server
	err = client.Connect(host, port)
	assertNoErr(err, t)
	assert(client.Connected(), t, "Client should be connected")
	time.Sleep(waitTime)

	// Check client address info
	host, port, err = client.GetAddr()
	assertNoErr(err, t)
	assert(client.sock.LocalAddr().String() == host+":"+strconv.Itoa(int(port)), t, "Address strings don't match")
	fmt.Printf("Client address: %s:%d\n", host, port)

	// Check connect event was received
	clientConnectEvent := <-serverEvent
	assertEq(clientConnectEvent, ServerEvent[any]{
		EventType: ServerConnect,
		ClientID:  0,
	}, t)

	// Check that addresses match
	host1, port1, err := client.GetAddr()
	assertNoErr(err, t)
	host2, port2, err := server.GetClientAddr(0)
	assertNoErr(err, t)
	assert(host1 == host2, t, "Client hosts do not match")
	assert(port1 == port2, t, "Client ports do not match")
	host3, port3, err := client.GetServerAddr()
	assertNoErr(err, t)
	host4, port4, err := server.GetAddr()
	assertNoErr(err, t)
	assert(host3 == host4, t, "Server hosts do not match")
	assert(port3 == port4, t, "Server ports do not match")

	// Disconnect from server
	err = client.Disconnect()
	assertNoErr(err, t)
	time.Sleep(waitTime)

	// Check disconnect event was received
	clientDisconnectEvent := <-serverEvent
	assertEq(clientDisconnectEvent, ServerEvent[any]{
		EventType: ServerDisconnect,
		ClientID:  0,
	}, t)

	// Stop server
	err = server.Stop()
	assertNoErr(err, t)
	assert(!server.Serving(), t, "Server should not be serving")
	time.Sleep(waitTime)
}

// Test sending messages between server and client
func TestSend(t *testing.T) {
	// Create server
	server, serverEvent := NewServer[int, string]()
	assert(!server.Serving(), t, "Server should not be serving")

	// Start server
	err := server.Start("127.0.0.1", 0)
	assertNoErr(err, t)
	assert(server.Serving(), t, "Server should be serving")
	time.Sleep(waitTime)

	// Check server address info
	host, port, err := server.GetAddr()
	assertNoErr(err, t)
	assert(server.sock.Addr().String() == host+":"+strconv.Itoa(int(port)), t, "Address strings don't match")
	fmt.Printf("Server address: %s:%d\n", host, port)

	// Create client
	client, clientEvent := NewClient[string, int]()
	assert(!client.Connected(), t, "Client should not be connected")

	// Connect to server
	err = client.Connect(host, port)
	assertNoErr(err, t)
	assert(client.Connected(), t, "Client should be connected")
	time.Sleep(waitTime)

	// Check client address info
	host, port, err = client.GetAddr()
	assertNoErr(err, t)
	assert(client.sock.LocalAddr().String() == host+":"+strconv.Itoa(int(port)), t, "Address strings don't match")
	fmt.Printf("Client address: %s:%d\n", host, port)

	// Check connect event was received
	clientConnectEvent := <-serverEvent
	assertEq(clientConnectEvent, ServerEvent[string]{
		EventType: ServerConnect,
		ClientID:  0,
	}, t)

	// Send message to client
	messageFromServer := 29275
	err = server.Send(messageFromServer)
	assertNoErr(err, t)
	time.Sleep(waitTime)

	// Receive message from server
	clientReceiveEvent1 := <-clientEvent
	assertEq(clientReceiveEvent1, ClientEvent[int]{
		EventType: ClientReceive,
		Data:      messageFromServer,
	}, t)
	time.Sleep(waitTime)

	// Send message to server
	messageFromClient := "Hello, server!"
	err = client.Send(messageFromClient)
	assertNoErr(err, t)
	time.Sleep(waitTime)

	// Receive message from client
	serverReceiveEvent := <-serverEvent
	assertEq(serverReceiveEvent, ServerEvent[string]{
		EventType: ServerReceive,
		ClientID:  0,
		Data:      messageFromClient,
	}, t)
	time.Sleep(waitTime)

	// Send response to client
	err = server.Send(len(serverReceiveEvent.Data))
	assertNoErr(err, t)
	time.Sleep(waitTime)

	// Receive response from server
	clientReceiveEvent2 := <-clientEvent
	assertEq(clientReceiveEvent2, ClientEvent[int]{
		EventType: ClientReceive,
		Data:      len(messageFromClient),
	}, t)
	time.Sleep(waitTime)

	// Disconnect from server
	err = client.Disconnect()
	assertNoErr(err, t)
	time.Sleep(waitTime)

	// Check disconnect event was received
	clientDisconnectEvent := <-serverEvent
	assertEq(clientDisconnectEvent, ServerEvent[string]{
		EventType: ServerDisconnect,
		ClientID:  0,
	}, t)

	// Stop server
	err = server.Stop()
	assertNoErr(err, t)
	assert(!server.Serving(), t, "Server should not be serving")
	time.Sleep(waitTime)
}

// Test sending large random messages between server and client
func TestLargeSend(t *testing.T) {
	// Create server
	server, serverEvent := NewServer[[]byte, []byte]()
	assert(!server.Serving(), t, "Server should not be serving")

	// Start server
	err := server.Start("127.0.0.1", 0)
	assertNoErr(err, t)
	assert(server.Serving(), t, "Server should be serving")
	time.Sleep(waitTime)

	// Check server address info
	host, port, err := server.GetAddr()
	assertNoErr(err, t)
	assert(server.sock.Addr().String() == host+":"+strconv.Itoa(int(port)), t, "Address strings don't match")
	fmt.Printf("Server address: %s:%d\n", host, port)

	// Create client
	client, clientEvent := NewClient[[]byte, []byte]()
	assert(!client.Connected(), t, "Client should not be connected")

	// Connect to server
	err = client.Connect(host, port)
	assertNoErr(err, t)
	assert(client.Connected(), t, "Client should be connected")
	time.Sleep(waitTime)

	// Check client address info
	host, port, err = client.GetAddr()
	assertNoErr(err, t)
	assert(client.sock.LocalAddr().String() == host+":"+strconv.Itoa(int(port)), t, "Address strings don't match")
	fmt.Printf("Client address: %s:%d\n", host, port)

	// Check connect event was received
	clientConnectEvent := <-serverEvent
	assertEq(clientConnectEvent, ServerEvent[[]byte]{
		EventType: ServerConnect,
		ClientID:  0,
	}, t)

	// Generate large messages
	largeMessageFromServerLength := rand.Int() % 65536
	largeMessageFromServer := make([]byte, largeMessageFromServerLength)
	n, err := rand.Read(largeMessageFromServer)
	assertNoErr(err, t)
	assertEq(n, largeMessageFromServerLength, t)
	fmt.Printf("Generated large message from server (%d bytes)\n", largeMessageFromServerLength)
	largeMessageFromClientLength := rand.Int() % 32768
	largeMessageFromClient := make([]byte, largeMessageFromClientLength)
	n, err = rand.Read(largeMessageFromClient)
	assertNoErr(err, t)
	assertEq(n, largeMessageFromClientLength, t)
	fmt.Printf("Generated large message from client (%d bytes)\n", largeMessageFromClientLength)

	// Send large message to client
	err = server.Send(largeMessageFromServer)
	assertNoErr(err, t)
	time.Sleep(waitTime)

	// Receive large message from server
	dataFromServer := <-clientEvent
	assertEq(dataFromServer, ClientEvent[[]byte]{
		EventType: ClientReceive,
		Data:      largeMessageFromServer,
	}, t)

	// Send large message to server
	err = client.Send(largeMessageFromClient)
	assertNoErr(err, t)
	time.Sleep(waitTime)

	// Receive large message from client
	dataFromClient := <-serverEvent
	assertEq(dataFromClient, ServerEvent[[]byte]{
		EventType: ServerReceive,
		ClientID:  0,
		Data:      largeMessageFromClient,
	}, t)

	// Disconnect from server
	err = client.Disconnect()
	assertNoErr(err, t)
	time.Sleep(waitTime)

	// Check disconnect event was received
	clientDisconnectEvent := <-serverEvent
	assertEq(clientDisconnectEvent, ServerEvent[[]byte]{
		EventType: ServerDisconnect,
		ClientID:  0,
	}, t)

	// Stop server
	err = server.Stop()
	assertNoErr(err, t)
	assert(!server.Serving(), t, "Server should not be serving")
	time.Sleep(waitTime)
}

// Test sending numerous messages
func TestSendingNumerousMessages(t *testing.T) {
	// Create server
	server, serverEvent := NewServer[int, int]()
	assert(!server.Serving(), t, "Server should not be serving")

	// Start server
	err := server.Start("127.0.0.1", 0)
	assertNoErr(err, t)
	assert(server.Serving(), t, "Server should be serving")
	time.Sleep(waitTime)

	// Check server address info
	host, port, err := server.GetAddr()
	assertNoErr(err, t)
	assert(server.sock.Addr().String() == host+":"+strconv.Itoa(int(port)), t, "Address strings don't match")
	fmt.Printf("Server address: %s:%d\n", host, port)

	// Create client
	client, clientEvent := NewClient[int, int]()
	assert(!client.Connected(), t, "Client should not be connected")

	// Connect to server
	err = client.Connect(host, port)
	assertNoErr(err, t)
	assert(client.Connected(), t, "Client should be connected")
	time.Sleep(waitTime)

	// Check client address info
	host, port, err = client.GetAddr()
	assertNoErr(err, t)
	assert(client.sock.LocalAddr().String() == host+":"+strconv.Itoa(int(port)), t, "Address strings don't match")
	fmt.Printf("Client address: %s:%d\n", host, port)

	// Check connect event was received
	clientConnectEvent := <-serverEvent
	assertEq(clientConnectEvent, ServerEvent[int]{
		EventType: ServerConnect,
		ClientID:  0,
	}, t)

	// Generate messages
	numServerMessages := (rand.Int() % 64) + 64
	numClientMessages := (rand.Int() % 128) + 128
	serverMessages := make([]int, numServerMessages)
	clientMessages := make([]int, numClientMessages)
	for i := 0; i < numServerMessages; i++ {
		serverMessages = append(serverMessages, rand.Int()%1024)
	}
	for i := 0; i < numClientMessages; i++ {
		clientMessages = append(clientMessages, rand.Int()%1024)
	}
	fmt.Printf("Generated %d server messages\n", numServerMessages)
	fmt.Printf("Generated %d client messages\n", numClientMessages)

	// Send messages
	for _, serverMessage := range serverMessages {
		err := client.Send(serverMessage)
		assertNoErr(err, t)
	}
	for _, clientMessage := range clientMessages {
		err := server.Send(clientMessage)
		assertNoErr(err, t)
	}
	time.Sleep(waitTime)

	// Receive messages from client
	for _, serverMessage := range serverMessages {
		serverReceiveEvent := <-serverEvent
		assertEq(serverReceiveEvent, ServerEvent[int]{
			EventType: ServerReceive,
			ClientID:  0,
			Data:      serverMessage,
		}, t)
	}

	// Receive messages from server
	for _, clientMessage := range clientMessages {
		clientReceiveEvent := <-clientEvent
		assertEq(clientReceiveEvent, ClientEvent[int]{
			EventType: ClientReceive,
			Data:      clientMessage,
		}, t)
	}

	// Disconnect from server
	err = client.Disconnect()
	assertNoErr(err, t)
	time.Sleep(waitTime)

	// Check disconnect event was received
	clientDisconnectEvent := <-serverEvent
	assertEq(clientDisconnectEvent, ServerEvent[int]{
		EventType: ServerDisconnect,
		ClientID:  0,
	}, t)

	// Stop server
	err = server.Stop()
	assertNoErr(err, t)
	assert(!server.Serving(), t, "Server should not be serving")
	time.Sleep(waitTime)
}

// Test sending custom types
func TestSendingCustomTypes(t *testing.T) {
	// Create server
	server, serverEvent := NewServer[custom, custom]()
	assert(!server.Serving(), t, "Server should not be serving")

	// Start server
	err := server.Start("127.0.0.1", 0)
	assertNoErr(err, t)
	assert(server.Serving(), t, "Server should be serving")
	time.Sleep(waitTime)

	// Check server address info
	host, port, err := server.GetAddr()
	assertNoErr(err, t)
	assert(server.sock.Addr().String() == host+":"+strconv.Itoa(int(port)), t, "Address strings don't match")
	fmt.Printf("Server address: %s:%d\n", host, port)

	// Create client
	client, clientEvent := NewClient[custom, custom]()
	assert(!client.Connected(), t, "Client should not be connected")

	// Connect to server
	err = client.Connect(host, port)
	assertNoErr(err, t)
	assert(client.Connected(), t, "Client should be connected")
	time.Sleep(waitTime)

	// Check client address info
	host, port, err = client.GetAddr()
	assertNoErr(err, t)
	assert(client.sock.LocalAddr().String() == host+":"+strconv.Itoa(int(port)), t, "Address strings don't match")
	fmt.Printf("Client address: %s:%d\n", host, port)

	// Check connect event was received
	clientConnectEvent := <-serverEvent
	assertEq(clientConnectEvent, ServerEvent[custom]{
		EventType: ServerConnect,
		ClientID:  0,
	}, t)

	// Messages
	serverMessage := custom{
		A: 123,
		B: "Hello, custom server class!",
		C: []string{"first server item", "second server item"},
	}
	clientMessage := custom{
		A: 456,
		B: "Hello, custom client class!",
		C: []string{"#1 client item", "client item #2", "(3) client item"},
	}

	// Send message to client
	err = server.Send(clientMessage)
	assertNoErr(err, t)
	time.Sleep(waitTime)

	// Receive message from server
	clientReceiveEvent1 := <-clientEvent
	assertEq(clientReceiveEvent1, ClientEvent[custom]{
		EventType: ClientReceive,
		Data:      clientMessage,
	}, t)
	time.Sleep(waitTime)

	// Send message to server
	err = client.Send(serverMessage)
	assertNoErr(err, t)
	time.Sleep(waitTime)

	// Receive message from client
	serverReceiveEvent := <-serverEvent
	assertEq(serverReceiveEvent, ServerEvent[custom]{
		EventType: ServerReceive,
		ClientID:  0,
		Data:      serverMessage,
	}, t)
	time.Sleep(waitTime)

	// Disconnect from server
	err = client.Disconnect()
	assertNoErr(err, t)
	time.Sleep(waitTime)

	// Check disconnect event was received
	clientDisconnectEvent := <-serverEvent
	assertEq(clientDisconnectEvent, ServerEvent[custom]{
		EventType: ServerDisconnect,
		ClientID:  0,
	}, t)

	// Stop server
	err = server.Stop()
	assertNoErr(err, t)
	assert(!server.Serving(), t, "Server should not be serving")
	time.Sleep(waitTime)
}

// Test having multiple clients connected, and process events from them individually
func TestMultipleClients(t *testing.T) {
	// Create server
	server, serverEvent := NewServer[int, string]()
	assert(!server.Serving(), t, "Server should not be serving")

	// Start server
	err := server.Start("127.0.0.1", 0)
	assertNoErr(err, t)
	assert(server.Serving(), t, "Server should be serving")
	time.Sleep(waitTime)

	// Check server address info
	host, port, err := server.GetAddr()
	assertNoErr(err, t)
	assert(server.sock.Addr().String() == host+":"+strconv.Itoa(int(port)), t, "Address strings don't match")
	fmt.Printf("Server address: %s:%d\n", host, port)

	// Create client
	client1, clientEvent1 := NewClient[string, int]()
	assert(!client1.Connected(), t, "Client should not be connected")

	// Connect to server
	err = client1.Connect(host, port)
	assertNoErr(err, t)
	assert(client1.Connected(), t, "Client should be connected")
	time.Sleep(waitTime)

	// Check client address info
	clientHost1, clientPort1, err := client1.GetAddr()
	assertNoErr(err, t)
	assert(client1.sock.LocalAddr().String() == clientHost1+":"+strconv.Itoa(int(clientPort1)), t, "Address strings don't match")
	fmt.Printf("Client address: %s:%d\n", clientHost1, clientPort1)

	// Check connect event was received
	clientConnectEvent1 := <-serverEvent
	assertEq(clientConnectEvent1, ServerEvent[string]{
		EventType: ServerConnect,
		ClientID:  0,
	}, t)

	// Check that first client addresses match
	host1, port1, err := client1.GetAddr()
	assertNoErr(err, t)
	host2, port2, err := server.GetClientAddr(0)
	assertNoErr(err, t)
	assert(host1 == host2, t, "Client 1 hosts do not match")
	assert(port1 == port2, t, "Client 1 ports do not match")
	host3, port3, err := client1.GetServerAddr()
	assertNoErr(err, t)
	host4, port4, err := server.GetAddr()
	assertNoErr(err, t)
	assert(host3 == host4, t, "Server hosts do not match")
	assert(port3 == port4, t, "Server ports do not match")

	// Create client
	client2, clientEvent2 := NewClient[string, int]()
	assert(!client2.Connected(), t, "Client should not be connected")

	// Connect to server
	err = client2.Connect(host, port)
	assertNoErr(err, t)
	assert(client2.Connected(), t, "Client should be connected")
	time.Sleep(waitTime)

	// Check client address info
	clientHost2, clientPort2, err := client2.GetAddr()
	assertNoErr(err, t)
	assert(client2.sock.LocalAddr().String() == clientHost2+":"+strconv.Itoa(int(clientPort2)), t, "Address strings don't match")
	fmt.Printf("Client address: %s:%d\n", clientHost2, clientPort2)

	// Check connect event was received
	clientConnectEvent2 := <-serverEvent
	assertEq(clientConnectEvent2, ServerEvent[string]{
		EventType: ServerConnect,
		ClientID:  1,
	}, t)

	// Check that second client addresses match
	host5, port5, err := client2.GetAddr()
	assertNoErr(err, t)
	host6, port6, err := server.GetClientAddr(1)
	assertNoErr(err, t)
	assert(host5 == host6, t, "Client 2 hosts do not match")
	assert(port5 == port6, t, "Client 2 ports do not match")
	host7, port7, err := client2.GetServerAddr()
	assertNoErr(err, t)
	host8, port8, err := server.GetAddr()
	assertNoErr(err, t)
	assert(host7 == host8, t, "Server hosts do not match")
	assert(port7 == port8, t, "Server ports do not match")

	// Send message from client 1
	messageFromClient1 := "Hello from client 1"
	err = client1.Send(messageFromClient1)
	assertNoErr(err, t)
	time.Sleep(waitTime)

	// Receive message from client 1
	serverMessageFromClient1 := <-serverEvent
	assertEq(serverMessageFromClient1, ServerEvent[string]{
		EventType: ServerReceive,
		ClientID:  0,
		Data:      messageFromClient1,
	}, t)

	// Send response back to client 1
	err = server.Send(len(serverMessageFromClient1.Data), serverMessageFromClient1.ClientID)
	assertNoErr(err, t)
	time.Sleep(waitTime)

	// Client 1 receive response
	serverReplyEvent1 := <-clientEvent1
	assertEq(serverReplyEvent1, ClientEvent[int]{
		EventType: ClientReceive,
		Data:      len(messageFromClient1),
	}, t)

	// Send message from client 2
	messageFromClient2 := "Hello from client 2"
	err = client2.Send(messageFromClient2)
	assertNoErr(err, t)
	time.Sleep(waitTime)

	// Receive message from client 2
	serverMessageFromClient2 := <-serverEvent
	assertEq(serverMessageFromClient2, ServerEvent[string]{
		EventType: ServerReceive,
		ClientID:  1,
		Data:      messageFromClient2,
	}, t)

	// Send response back to client 2
	err = server.Send(len(serverMessageFromClient2.Data), serverMessageFromClient2.ClientID)
	assertNoErr(err, t)
	time.Sleep(waitTime)

	// Client 2 receive response
	serverReplyEvent2 := <-clientEvent2
	assertEq(serverReplyEvent2, ClientEvent[int]{
		EventType: ClientReceive,
		Data:      len(messageFromClient2),
	}, t)

	// Send message to all clients
	messageFromServer := 29275
	err = server.Send(messageFromServer)
	assertNoErr(err, t)
	time.Sleep(waitTime)

	// Client 1 receive message
	serverMessage1 := <-clientEvent1
	assertEq(serverMessage1, ClientEvent[int]{
		EventType: ClientReceive,
		Data:      messageFromServer,
	}, t)

	// Client 2 receive message
	serverMessage2 := <-clientEvent2
	assertEq(serverMessage2, ClientEvent[int]{
		EventType: ClientReceive,
		Data:      messageFromServer,
	}, t)

	// Client 1 disconnect from server
	err = client1.Disconnect()
	assertNoErr(err, t)
	time.Sleep(waitTime)

	// Check client 1 disconnect event was received
	clientDisconnectEvent1 := <-serverEvent
	assertEq(clientDisconnectEvent1, ServerEvent[string]{
		EventType: ServerDisconnect,
		ClientID:  0,
	}, t)

	// Client 2 disconnect from server
	err = client2.Disconnect()
	assertNoErr(err, t)
	time.Sleep(waitTime)

	// Check client 2 disconnect event was received
	clientDisconnectEvent2 := <-serverEvent
	assertEq(clientDisconnectEvent2, ServerEvent[string]{
		EventType: ServerDisconnect,
		ClientID:  1,
	}, t)

	// Stop server
	err = server.Stop()
	assertNoErr(err, t)
	assert(!server.Serving(), t, "Server should not be serving")
	time.Sleep(waitTime)
}

// Test removing a client from the server
func TestRemoveClient(t *testing.T) {
	// Create server
	server, serverEvent := NewServer[any, any]()
	assert(!server.Serving(), t, "Server should not be serving")

	// Start server
	err := server.Start("127.0.0.1", 0)
	assertNoErr(err, t)
	assert(server.Serving(), t, "Server should be serving")
	time.Sleep(waitTime)

	// Check server address info
	host, port, err := server.GetAddr()
	assertNoErr(err, t)
	assert(server.sock.Addr().String() == host+":"+strconv.Itoa(int(port)), t, "Address strings don't match")
	fmt.Printf("Server address: %s:%d\n", host, port)

	// Create client
	client, clientEvent := NewClient[any, any]()
	assert(!client.Connected(), t, "Client should not be connected")

	// Connect to server
	err = client.Connect(host, port)
	assertNoErr(err, t)
	assert(client.Connected(), t, "Client should be connected")
	time.Sleep(waitTime)

	// Check client address info
	host, port, err = client.GetAddr()
	assertNoErr(err, t)
	assert(client.sock.LocalAddr().String() == host+":"+strconv.Itoa(int(port)), t, "Address strings don't match")
	fmt.Printf("Client address: %s:%d\n", host, port)

	// Check connect event was received
	clientConnectEvent := <-serverEvent
	assertEq(clientConnectEvent, ServerEvent[any]{
		EventType: ServerConnect,
		ClientID:  0,
	}, t)

	// Remove the client
	err = server.RemoveClient(0)
	assertNoErr(err, t)
	time.Sleep(waitTime)

	// Receive the client disconnected event
	disconnectedEvent := <-clientEvent
	assertEq(disconnectedEvent, ClientEvent[any]{
		EventType: ClientDisconnected,
	}, t)

	// Receive the server disconnect event
	disconnectEvent := <-serverEvent
	assertEq(disconnectEvent, ServerEvent[any]{
		EventType: ServerDisconnect,
		ClientID:  0,
	}, t)

	// Check the client is not connected
	assert(!client.Connected(), t, "Client should not be connected")

	// Stop server
	err = server.Stop()
	assertNoErr(err, t)
	assert(!server.Serving(), t, "Server should not be serving")
	time.Sleep(waitTime)
}

// Test stopping a server while a client is connected
func TestStopServerWhileClientConnected(t *testing.T) {
	// Create server
	server, serverEvent := NewServer[any, any]()
	assert(!server.Serving(), t, "Server should not be serving")

	// Start server
	err := server.Start("127.0.0.1", 0)
	assertNoErr(err, t)
	assert(server.Serving(), t, "Server should be serving")
	time.Sleep(waitTime)

	// Check server address info
	host, port, err := server.GetAddr()
	assertNoErr(err, t)
	assert(server.sock.Addr().String() == host+":"+strconv.Itoa(int(port)), t, "Address strings don't match")
	fmt.Printf("Server address: %s:%d\n", host, port)

	// Create client
	client, clientEvent := NewClient[any, any]()
	assert(!client.Connected(), t, "Client should not be connected")

	// Connect to server
	err = client.Connect(host, port)
	assertNoErr(err, t)
	assert(client.Connected(), t, "Client should be connected")
	time.Sleep(waitTime)

	// Check client address info
	host, port, err = client.GetAddr()
	assertNoErr(err, t)
	assert(client.sock.LocalAddr().String() == host+":"+strconv.Itoa(int(port)), t, "Address strings don't match")
	fmt.Printf("Client address: %s:%d\n", host, port)

	// Check connect event was received
	clientConnectEvent := <-serverEvent
	assertEq(clientConnectEvent, ServerEvent[any]{
		EventType: ServerConnect,
		ClientID:  0,
	}, t)

	// Stop server
	err = server.Stop()
	assertNoErr(err, t)
	assert(!server.Serving(), t, "Server should not be serving")
	time.Sleep(waitTime)

	// Check disconnected event was received
	disconnectedEvent := <-clientEvent
	assertEq(disconnectedEvent, ClientEvent[any]{
		EventType: ClientDisconnected,
	}, t)
	assert(!client.Connected(), t, "Client should not be connected")
}
