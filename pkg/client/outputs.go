package client

import (
	"context"
	"math/rand"
	"strconv"

	"github.com/falcosecurity/client-go/pkg/api/outputs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// OutputsClient is a client wrapping the gRPC outputs service client.
type OutputsClient struct {
	c         outputs.ServiceClient
	sessionID string
}

// Outputs requests Falco outputs to the Falco gRPC Outputs API.
func (o *OutputsClient) Outputs(ctx context.Context, req *outputs.Request, opts ...grpc.CallOption) (outputs.Service_OutputsClient, error) {
	reqID := strconv.Itoa(rand.Intn(1000) + 1)
	ctx = metadata.AppendToOutgoingContext(ctx, "session_id", o.sessionID, "request_id", reqID)
	return o.c.Outputs(ctx, req, opts...)
}
