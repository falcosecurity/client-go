package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"

	"github.com/falcosecurity/client-go/pkg/api/outputs"
	"github.com/falcosecurity/client-go/pkg/api/version"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client is a wrapper for the gRPC connection
// it allows to connect to a Falco gRPC server.
// It is created using the function NewForConfig(config *Config) .
type Client struct {
	conn *grpc.ClientConn
}

// Config is the configuration definition for connecting to a Falco gRPC server.
type Config struct {
	Hostname   string
	Port       uint16
	CertFile   string
	KeyFile    string
	CARootFile string
}

const targetFormat = "%s:%d"

// NewForConfig is used to create a new Falco gRPC client.
func NewForConfig(config *Config) (*Client, error) {
	rand.Seed(time.Now().UnixNano())

	certificate, err := tls.LoadX509KeyPair(
		config.CertFile,
		config.KeyFile,
	)
	if err != nil {
		return nil, fmt.Errorf("error loading the X.509 key pair: %v", err)
	}

	certPool := x509.NewCertPool()
	rootCA, err := ioutil.ReadFile(config.CARootFile)
	if err != nil {
		return nil, fmt.Errorf("error reading the CA Root file certificate: %v", err)
	}

	ok := certPool.AppendCertsFromPEM(rootCA)
	if !ok {
		return nil, fmt.Errorf("error appending the root CA to the certificate pool")
	}

	transportCreds := credentials.NewTLS(&tls.Config{
		ServerName:   config.Hostname,
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})

	dialOption := grpc.WithTransportCredentials(transportCreds)
	conn, err := grpc.Dial(fmt.Sprintf(targetFormat, config.Hostname, config.Port), dialOption)
	if err != nil {
		return nil, fmt.Errorf("error dialing server: %v", err)
	}

	return &Client{
		conn,
	}, nil
}

// Outputs is the client for the Falco gRPC Outputs API.
// When using it you can use `outputs()` to receive a stream of Falco output events.
func (c *Client) Outputs() (*OutputsClient, error) {
	if err := c.checkConn(); err != nil {
		return nil, err
	}
	oc := &OutputsClient{
		c: outputs.NewServiceClient(c.conn),
		sessionID: strconv.Itoa(rand.Intn(1000) + 1),
	}

	return oc, nil
}

// Version it the client for the Falco gRPC Version API.
// When using it you can use `version()` to receive the Falco version.
func (c *Client) Version() (*VersionClient, error) {
	if err := c.checkConn(); err != nil {
		return nil, err
	}
	vc := &VersionClient{
		c:         version.NewServiceClient(c.conn),
		sessionID: strconv.Itoa(rand.Intn(1000) + 1),
	}

	return vc, nil
}

// Close the connection to the Falco gRPC server.
func (c *Client) Close() error {
	if err := c.checkConn(); err != nil {
		return err
	}
	return c.conn.Close()
}

func (c *Client) checkConn() error {
	if c.conn == nil {
		return fmt.Errorf("missing connection for the current client")
	}
	return nil
}
