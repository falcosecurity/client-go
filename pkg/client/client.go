package client

import (
	"fmt"

	"github.com/falcosecurity/client-go/pkg/api/output"
	"google.golang.org/grpc"
)

// Client is a wrapper for the gRPC connection
// it allows to connect to a Falco gRPC server.
// It is created using the function NewForConfig(config *Config) .
type Client struct {
	conn *grpc.ClientConn
}

// Config is the configuration definition for connecting to a Falco gRPC server.
type Config struct {
	Target  string
	Options []grpc.DialOption
}

// NewForConfig is used to create a new Falco gRPC client.
func NewForConfig(config *Config) (*Client, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(config.Target, config.Options...)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn,
	}, nil
}

// Output is the client for Falco Outputs.
// When using it you can use `subscribe()` to receive a stream of falco output events.
func (c *Client) Output() (output.FalcoOutputServiceClient, error) {
	if err := c.checkConn(); err != nil {
		return nil, err
	}
	return output.NewFalcoOutputServiceClient(c.conn), nil
}

// Close the connection to the falco gRPC server.
func (c *Client) Close() error {
	if err := c.checkConn(); err != nil {
		return err
	}
	return c.conn.Close()
}

func (c *Client) checkConn() error {
	if c.conn == nil {
		return fmt.Errorf("missing gRPC connection for the current client")
	}
	return nil
}
