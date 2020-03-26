package godtp

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
)

const lenSize = 5
const msgLabel = "message"

// Encode an object so it can be sent through a socket
func Encode(object interface{}) ([]byte, error) {
	return json.Marshal(object)
}

// Decode an object that came through a socket
func Decode(bytestring []byte, object interface{}) error {
	return json.Unmarshal(bytestring, object)
}

// Convert decimal to ASCII
func decToASCII(dec uint64) []byte {
	ascii := make([]byte, lenSize)
	for i := lenSize - 1; i >= 0; i-- {
		ascii[i] = uint8(dec)
		dec >>= 8
	}
	return ascii
}

// Convert ASCII to decimal
func asciiToDec(ascii []byte) uint64 {
	var dec uint64 = 0
	for i := 0; i < lenSize; i++ {
		dec <<= 8
		dec += uint64(ascii[i])
	}
	return dec
}

// Parse an address
func parseAddr(addr string) (string, uint16, error) {
	index := strings.LastIndex(addr, ":")
	if index > -1 {
		port, err := strconv.Atoi(addr[index + 1:])
		if err == nil {
			return addr[:index], uint16(port), nil
		}
		return "", 0, fmt.Errorf("Port conversion failed")
	}
	return "", 0, fmt.Errorf("No port found")
}

// Generate a new 2048-bit RSA key pair
func newKeys() (*rsa.PrivateKey, error) {
	rng := rand.Reader
	priv, err := rsa.GenerateKey(rng, 2048)
	return priv, err
}

// Encrypt using RSA
func asymmetricEncrypt(pub rsa.PublicKey, plaintext []byte) ([]byte, error) {
	rng := rand.Reader
	label := []byte(msgLabel)
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &pub, plaintext, label)
	return ciphertext, err
}

// Decrypt using RSA
func asymmetricDecrypt(priv *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	rng := rand.Reader
	label := []byte(msgLabel)
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, priv, ciphertext, label)
	return plaintext, err
}

// Generate a new 256-bit AES key
func newKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	return key, err
}

// Encrypt using AES
func encrypt(key, plaintext []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return []byte{}, err
	}

    ciphertext := make([]byte, aes.BlockSize + len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return []byte{}, err
	}

    cfb := cipher.NewCFBEncrypter(block, iv)
    cfb.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
    return ciphertext, nil
}

// Decrypt using AES
func decrypt(key []byte, ciphertext []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return []byte{}, err
	}

    if len(ciphertext) < aes.BlockSize {
        return []byte{}, fmt.Errorf("ciphertext too short")
	}

    iv := ciphertext[:aes.BlockSize]
    plaintext := ciphertext[aes.BlockSize:]
    cfb := cipher.NewCFBDecrypter(block, iv)
    cfb.XORKeyStream(plaintext, plaintext)
    return plaintext, nil
}

// Close a byte slice channel
func closeBytesChan(ch chan []byte) bool {
	open := make(chan bool)
	var wg sync.WaitGroup

	go func() {
		wg.Add(1)

		defer func() {
			open <- recover() == nil
			wg.Done()
		}()

		close(ch)
	}()

	wg.Wait()
	return <-open
}

// Close a byte slice channel channel
func closeBytesChanChan(ch chan chan []byte) bool {
	open := make(chan bool)
	var wg sync.WaitGroup

	go func() {
		wg.Add(1)

		defer func() {
			open <- recover() == nil
			wg.Done()
		}()

		close(ch)
	}()

	wg.Wait()
	return <-open
}
