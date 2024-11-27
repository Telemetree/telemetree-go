package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/TONSolutions/telemetree-go/telemetree/enity"
)

const eventSource = "Go SDK"

// Payload represents the tracking event data
type Payload struct {
	ApplicationID string  `json:"application_id"`
	Datetime      int64   `json:"datetime"`
	SessionID     int     `json:"session_id"`
	TelegramID    int     `json:"telegram_id"`
	Referrer      *int    `json:"referrer,omitempty"`
	EventSource   string  `json:"event_source"`
	IsPremium     bool    `json:"is_premium"`
	EventType     string  `json:"event_type"`
	Username      *string `json:"username,omitempty"`
	Firstname     *string `json:"firstname,omitempty"`
	Lastname      *string `json:"lastname,omitempty"`
	Language      *string `json:"language,omitempty"`
	ReferrerType  *string `json:"referrer_type,omitempty"`
}

// EncryptedPayload represents the final encrypted data structure
type EncryptedPayload struct {
	EncryptedData []byte `json:"body"`
	EncryptedKey  []byte `json:"key"`
	EncryptedIV   []byte `json:"iv"`
}

func buildPayload(event enity.Event, applicationID string) Payload {
	payload := Payload{
		ApplicationID: applicationID,
		Datetime:      time.Now().UTC().Unix(),
		EventSource:   eventSource,
		SessionID:     getCurrentSessionID(),
		IsPremium:     event.IsPremium,
		TelegramID:    event.TelegramID,
		EventType:     event.EventType,
	}

	if event.Username != "" {
		payload.Username = &event.Username
	}
	if event.Firstname != "" {
		payload.Firstname = &event.Firstname
	}
	if event.Lastname != "" {
		payload.Lastname = &event.Lastname
	}

	if event.Language != "" {
		payload.Language = &event.Language
	}

	if event.ReferrerType != "" {
		payload.ReferrerType = &event.ReferrerType
	}

	if event.Referrer != "" {
		referrer, _ := strconv.Atoi(event.Referrer)
		payload.Referrer = &referrer
	}

	return payload
}

func getCurrentSessionID() int {
	return int(time.Now().UTC().UnixNano() / int64(time.Millisecond))
}

// generateAESKeyAndIV generates a random 16-byte AES key and IV
func generateAESKeyAndIV() (key []byte, iv []byte, err error) {
	key = make([]byte, 16) // Using 16 bytes to match JS implementation
	if _, err = rand.Read(key); err != nil {
		return nil, nil, err
	}

	iv = make([]byte, aes.BlockSize)
	if _, err = rand.Read(iv); err != nil {
		return nil, nil, err
	}

	return key, iv, nil
}

// encryptWithAES encrypts a message using AES in CBC mode
func encryptWithAES(key []byte, iv []byte, message string) ([]byte, error) {
	// Create padded message
	messageBytes := []byte(message)
	blockSize := aes.BlockSize
	padding := blockSize - len(messageBytes)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	messageBytes = append(messageBytes, padText...)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Encrypt using CBC mode
	encryptedData := make([]byte, len(messageBytes))
	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(encryptedData, messageBytes)

	return encryptedData, nil
}

func encryptWithRSA(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	// Convert bytes to hex string first, matching JS implementation
	hexStr := fmt.Sprintf("%x", data)

	// Encrypt the hex string
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(hexStr))
}

func parseRSAPublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
	publicKeyPEM = strings.ReplaceAll(publicKeyPEM, "\\n", "\n")

	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}

	if block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("invalid key type: expected RSA PUBLIC KEY")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKCS1 public key: %w", err)
	}
	return publicKey, nil
}

func PrepareEncryptedPayload(publicKeyPEM string, applicationID string, event enity.Event) ([]byte, error) {
	publicKey, err := parseRSAPublicKey(strings.ReplaceAll(publicKeyPEM, "\\n", "\n"))
	payload := buildPayload(event, applicationID)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	aesKey, iv, err := generateAESKeyAndIV()
	if err != nil {
		return nil, err
	}

	encryptedData, err := encryptWithAES(aesKey, iv, string(payloadBytes))
	if err != nil {
		return nil, err
	}

	encryptedKey, err := encryptWithRSA(publicKey, aesKey)
	if err != nil {
		return nil, err
	}

	encryptedIV, err := encryptWithRSA(publicKey, iv)
	if err != nil {
		return nil, err
	}

	jsonResult, err := json.Marshal(EncryptedPayload{
		EncryptedData: encryptedData,
		EncryptedKey:  encryptedKey,
		EncryptedIV:   encryptedIV,
	})
	if err != nil {
		return nil, err
	}

	return jsonResult, nil
}
