package mx


import (
    "crypto/md5"
    "encoding/binary"
    "encoding/hex"
)

// GetMxToken returns a hex-encoded MD5 of the given seed, truncated to 'length' characters.
// seed:   e.g. account ID, timestamp, or random int.
// length: desired token length (must be <= 32, since MD5 hex is 32 chars).
func GetMxToken(seed int64, length int) string {
    const hexLen = md5.Size * 2 // 16 bytes * 2 = 32 hex chars
    if length <= 0 || length > hexLen {
        return ""
    }

    // Convert seed to 8-byte big-endian
    buf := make([]byte, 8)
    binary.BigEndian.PutUint64(buf, uint64(seed))

    // Compute MD5
    sum := md5.Sum(buf)           // [16]byte
    hexStr := hex.EncodeToString(sum[:]) // 32-char string

    // Truncate
    return hexStr[:length]
}

var Key string
var Docker string
type ProtoMessage interface {
	String() string
	SetPacket(packet Message)
}

type Message = ProtoMessage
