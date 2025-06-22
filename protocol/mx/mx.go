package mx


import (
	"crypto/rand"
	"encoding/hex"
)

const maxTokenLength = 64

// GetMxToken returns a random hex string of the given length (max 64).
// Ignores the seed for now.
func GetMxToken(_ int64, length int) string {
	if length <= 0 {
		length = 32 // fallback default
	} else if length > maxTokenLength {
		length = maxTokenLength
	}

	// Generate random bytes (1 hex char = 4 bits → 2 chars = 1 byte)
	byteLen := (length + 1) / 2
	buf := make([]byte, byteLen)

	if _, err := rand.Read(buf); err != nil {
		return "" // fallback on failure
	}

	token := hex.EncodeToString(buf)
	return token[:length]
}

type ProtoMessage interface {
	String() string
	SetPacket(packet Message)
}

type Message = ProtoMessage
