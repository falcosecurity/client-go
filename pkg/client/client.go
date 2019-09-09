package client

import (
	"fmt"

	"github.com/falcosecurity/client-go/pkg/api/output"
	"google.golang.org/grpc"
)

// Client ...
type Client struct {
	conn *grpc.ClientConn
}

// Config ...
type Config struct {
	Target  string
	Options []grpc.DialOption
}

// NewForConfig ...
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

// Output ...
func (c *Client) Output() (output.FalcoOutputServiceClient, error) {
	if err := c.checkConn(); err != nil {
		return nil, err
	}
	return output.NewFalcoOutputServiceClient(c.conn), nil
}

// Close ...
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
