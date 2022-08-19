package godtp

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"io"
)

// Message label
const msgLabel = "message"

// Generate a new 2048-bit RSA key pair
func newRSAKeys() (*rsa.PrivateKey, error) {
	rng := rand.Reader
	privateKey, err := rsa.GenerateKey(rng, 2048)
	return privateKey, err
}

// Encrypt using RSA
func rsaEncrypt(publicKey rsa.PublicKey, plaintext []byte) ([]byte, error) {
	rng := rand.Reader
	label := []byte(msgLabel)
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &publicKey, plaintext, label)
	return ciphertext, err
}

// Decrypt using RSA
func rsaDecrypt(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	rng := rand.Reader
	label := []byte(msgLabel)
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, privateKey, ciphertext, label)
	return plaintext, err
}

// Generate a new 256-bit AES key
func newAESKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	return key, err
}

// Encrypt using AES
func aesEncrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return []byte{}, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext, nil
}

// Decrypt using AES
func aesDecrypt(key []byte, ciphertext []byte) ([]byte, error) {
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
