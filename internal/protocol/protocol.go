package protocol

import (
	"encoding/base64"
	"fmt"

	"github.com/elizandrodantas/machine-controller-v2/internal/util"
)

var (
	// version is the version of the protocol.
	version = []byte("2")

	// keyLen is the length of the key.
	_keyLen = 16
	_ivLen  = 16
)

type protocol struct {
	key     []byte
	iv      []byte
	message []byte
	version []byte
}

// IProtocol represents a protocol interface.
type IProtocol interface {
	GetKey() []byte
	GetIV() []byte
	GetMessage() []byte
	GetVersion() []byte
	GetVersionString() string
	ToByteArray() []byte
	ToBase64() string
	IsVersionValid() bool
}

// ImportProtocolBase64 decodes a base64-encoded string and imports the protocol.
// It returns the imported protocol and any error encountered during decoding or importing.
func ImportProtocolBase64(data string) (IProtocol, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("error decoding base64: %v", err)
	}

	return ImportProtocol(decoded)
}

// ImportProtocol is a function that imports a protocol from the given data.
func ImportProtocol(data []byte) (IProtocol, error) {
	if len(data) < (_keyLen + _ivLen + len(version)) {
		return nil, fmt.Errorf("protocol invalid, data too short to contain key and iv")
	}

	protocol := &protocol{
		version: data[:1],
		key:     data[1 : _keyLen+1],
		iv:      data[len(data)-_ivLen:],
		message: data[1+_keyLen : len(data)-_ivLen],
	}

	return protocol, nil
}

// CreateProtocol creates a new protocol instance with the given key, iv, and message.
// If the key is not provided, a random key will be generated.
// If the iv is not provided, a random iv will be generated.
// An error is returned if the key or iv length is invalid.
func CreateProtocol(key, iv, message []byte) (IProtocol, error) {
	if key != nil && len(key) != 16 {
		return nil, fmt.Errorf("key must be 16 bytes")
	}

	if iv != nil && len(iv) != 16 {
		return nil, fmt.Errorf("iv must be 16 bytes")
	}

	if key == nil {
		key = util.RandomByte(_keyLen)
	}

	if iv == nil {
		iv = util.RandomByte(_ivLen)
	}

	return &protocol{
		version: version,
		key:     key,
		iv:      iv,
		message: message,
	}, nil
}

// GetKey returns the key associated with the protocol.
func (p *protocol) GetKey() []byte {
	return p.key
}

// GetIV returns the initialization vector (IV) used by the protocol.
func (p *protocol) GetIV() []byte {
	return p.iv
}

// GetMessage returns the message stored in the protocol.
func (p *protocol) GetMessage() []byte {
	return p.message
}

// GetVersion returns the version of the protocol.
func (p *protocol) GetVersion() []byte {
	return p.version
}

// GetVersionString returns the version string of the protocol.
func (p *protocol) GetVersionString() string {
	return string(p.version)
}

// IsVersionValid checks if the protocol version is valid.
// It returns true if the version is "2", otherwise it returns false.
func (p *protocol) IsVersionValid() bool {
	return string(p.version) == "2"
}

// ToByteArray converts the protocol struct to a byte array.
// It concatenates the version, key, message, and iv fields of the protocol struct
// and returns the resulting byte array.
func (p *protocol) ToByteArray() []byte {
	var output = make([]byte, 0)
	output = append(output, p.version...)
	output = append(output, p.key...)
	output = append(output, p.message...)
	output = append(output, p.iv...)

	return output
}

// ToBase64 encodes the protocol object to a base64 string.
func (p *protocol) ToBase64() string {
	return base64.StdEncoding.EncodeToString(p.ToByteArray())
}
