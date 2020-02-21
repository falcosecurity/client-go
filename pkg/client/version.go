package client

import (
	"context"
	"math/rand"
	"strconv"

	"github.com/falcosecurity/client-go/pkg/api/version"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// VersionClient is a client wrapping the gRPC version service client.
type VersionClient struct {
	c         version.ServiceClient
	sessionID string
}

// Version requests the Falco version to the Falco gRPC Version API.
func (v *VersionClient) Version(ctx context.Context, req *version.Request, opts ...grpc.CallOption) (*version.Response, error) {
	reqID := strconv.Itoa(rand.Intn(1000) + 1)
	ctx = metadata.AppendToOutgoingContext(ctx, "session_id", v.sessionID, "request_id", reqID)
	return v.c.Version(ctx, req, opts...)
}
