package client

import (
	"context"
	"math/rand"
	"strconv"

	"github.com/falcosecurity/client-go/pkg/api/inputs"
	"github.com/falcosecurity/client-go/pkg/api/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// InputsClient is a client wrapping the gRPC inputs service client.
type InputsClient struct {
	c         inputs.ServiceClient
	sessionID string
	source    schema.Source
}

// Input sends ...
func (i *InputsClient) Input(ctx context.Context, req *inputs.Request, opts ...grpc.CallOption) (*inputs.Response, error) {
	reqID := strconv.Itoa(rand.Intn(1000) + 1)
	ctx = metadata.AppendToOutgoingContext(ctx, "session_id", i.sessionID, "request_id", reqID)
	return i.c.Input(ctx, req, opts...)
}
