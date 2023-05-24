package client

import "errors"

type Config struct {
	// Endpoints are the endpoints of the servers.
	Endpoints []string
	// EnableEncrypted indicates whether to enable encrypted communication.
	EnableEncrypted bool
	// AesKey is the key for AES encryption. Only valid when EnableEncrypted is true.
	AesKey string
	// AesIv is the iv for AES encryption. Only valid when EnableEncrypted is true.
	AesIv string
	// ContentType is the content type of the request. Please select one of the following: protobuf, json; default is protobuf.
	ContentType string
}

var (
	ErrNoEndpoint = errors.New("no endpoint available")
)

func (c *Config) Validate() error {
	if len(c.Endpoints) == 0 {
		return ErrNoEndpoint
	}
	if c.ContentType != "json" {
		c.ContentType = "protobuf"
	}
	return nil
}
