package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/falcosecurity/client-go/pkg/api/outputs"
	"github.com/falcosecurity/client-go/pkg/api/version"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client is a wrapper for the gRPC connection
// it allows to connect to a Falco gRPC server.
// It is created using the function `NewForConfig(context.Context, *Config)`.
type Client struct {
	conn                 *grpc.ClientConn
	versionServiceClient version.ServiceClient
	outputsServiceClient outputs.ServiceClient
}

// Config is the configuration definition for connecting to a Falco gRPC server.
type Config struct {
	Hostname       string
	Port           uint16
	CertFile       string
	KeyFile        string
	CARootFile     string
	UnixSocketPath string
	DialOptions    []grpc.DialOption
}

const targetFormat = "%s:%d"

// NewForConfig is used to create a new Falco gRPC client.
func NewForConfig(ctx context.Context, config *Config) (*Client, error) {
	if len(config.UnixSocketPath) > 0 {
		return newUnixSocketClient(ctx, config)
	}
	return newNetworkClient(ctx, config)
}

func newUnixSocketClient(ctx context.Context, config *Config) (*Client, error) {
	dialOptions := append(config.DialOptions, grpc.WithInsecure())
	conn, err := grpc.DialContext(ctx, config.UnixSocketPath, dialOptions...)
	if err != nil {
		return nil, fmt.Errorf("error dialing server: %v", err)
	}
	return &Client{
		conn: conn,
	}, nil
}

func newNetworkClient(ctx context.Context, config *Config) (*Client, error) {
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

	dialOptions := append(config.DialOptions, grpc.WithTransportCredentials(transportCreds))
	conn, err := grpc.DialContext(ctx, fmt.Sprintf(targetFormat, config.Hostname, config.Port), dialOptions...)
	if err != nil {
		return nil, fmt.Errorf("error dialing server: %v", err)
	}

	return &Client{
		conn: conn,
	}, nil
}

// Outputs is the client for Falco Outputs.
// When using it you can use `Sub()` or `Get()` to receive a stream of Falco output events.
func (c *Client) Outputs() (outputs.ServiceClient, error) {
	if err := c.checkConn(); err != nil {
		return nil, err
	}
	if c.outputsServiceClient == nil {
		c.outputsServiceClient = outputs.NewServiceClient(c.conn)
	}
	return c.outputsServiceClient, nil
}

// OutputsWatchCallback is passed to OutputsWatch to perform an
// action for each *outputs.Response while retrieving a stream outputs
type OutputsWatchCallback func(res *outputs.Response) error

// OutputsWatch allows to watch and process a stream of *outputs.Response
// using a callback function of type OutputsWatchCallback.
//
// The timeout parameter specifies the frequency of the watch operation.
func (c *Client) OutputsWatch(ctx context.Context,
	cb OutputsWatchCallback,
	timeout time.Duration,
	opts ...grpc.CallOption) error {
	oc, err := c.Outputs()
	if err != nil {
		return err
	}
	fcs, err := oc.Sub(ctx, opts...)
	if err != nil {
		return err
	}
	resCh := make(chan *outputs.Response, 1)
	errCh := make(chan error, 1)

	go func() {
		for {
			res, err := fcs.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				errCh <- fmt.Errorf("error closing stream after EOF: %v", err)
			}
			resCh <- res
		}
	}()

	for {
		select {
		case res := <-resCh:
			err := cb(res)
			if err != nil {
				return err
			}
		case err := <-errCh:
			return err
		case <-time.After(timeout):
			fcs.Send(&outputs.Request{})
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Version it the client for Falco Version API.
// When using it you can use `version()` to receive the Falco version.
func (c *Client) Version() (version.ServiceClient, error) {
	if err := c.checkConn(); err != nil {
		return nil, err
	}
	if c.versionServiceClient == nil {
		c.versionServiceClient = version.NewServiceClient(c.conn)
	}
	return c.versionServiceClient, nil
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
		return fmt.Errorf("missing connection for the current client")
	}
	return nil
}
