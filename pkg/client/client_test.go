package client

import (
	"context"
	"testing"

	"github.com/falcosecurity/client-go/pkg/api/outputs"
	outputsmock "github.com/falcosecurity/client-go/pkg/api/outputs/mocks"
	"github.com/falcosecurity/client-go/pkg/api/version"
	versionmock "github.com/falcosecurity/client-go/pkg/api/version/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestClient_Version(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockVersionClient := versionmock.NewMockServiceClient(ctrl)
	mockVersionClient.EXPECT().Version(
		gomock.Any(),
		gomock.Any(),
	).Return(&version.Response{Version: "mockVersion"}, nil)

	c := Client{
		conn:                 &grpc.ClientConn{},
		versionServiceClient: mockVersionClient,
	}

	versionServiceClient, err := c.Version()
	assert.Nil(t, err)

	res, err := versionServiceClient.Version(context.Background(), &version.Request{})
	assert.Nil(t, err)
	assert.Equal(t, res.Version, "mockVersion", "version does not match mockVersion")
}

func TestClient_Outputs(t *testing.T) {
	ctrl := gomock.NewController(t)

	streamStub := outputsmock.NewMockService_GetClient(ctrl)
	streamStub.EXPECT().Recv().Return(&outputs.Response{Rule: "testRule", Output: "testOutput"}, nil)

	mockOutputClient := outputsmock.NewMockServiceClient(ctrl)
	mockOutputClient.EXPECT().Get(
		gomock.Any(),
		gomock.Any(),
	).Return(streamStub, nil)

	c := Client{
		conn:                 &grpc.ClientConn{},
		outputsServiceClient: mockOutputClient,
	}

	outputsServiceClient, err := c.Outputs()
	assert.Nil(t, err)

	stream, err := outputsServiceClient.Get(context.Background(), &outputs.Request{})
	assert.Nil(t, err)
	res, err := stream.Recv()
	assert.Nil(t, err)
	assert.Equal(t, res.Rule, "testRule", "rule does not match testRule")
	assert.Equal(t, res.Output, "testOutput", "output does not match testOutput")
}
