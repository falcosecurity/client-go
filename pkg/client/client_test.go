package client

import (
	"context"
	"testing"

	"github.com/falcosecurity/client-go/pkg/api/output"
	"github.com/falcosecurity/client-go/pkg/api/version"
	outputmock "github.com/falcosecurity/client-go/pkg/api/output/mocks"
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

func TestClient_Output(t *testing.T) {
	ctrl := gomock.NewController(t)

	streamStub := outputmock.NewMockService_SubscribeClient(ctrl)
	streamStub.EXPECT().Recv().Return(&output.Response{Rule: "testRule", Output: "testOutput"}, nil)

	mockOutputClient := outputmock.NewMockServiceClient(ctrl)
	mockOutputClient.EXPECT().Subscribe(
		gomock.Any(),
		gomock.Any(),
	).Return(streamStub, nil)

	c := Client{
		conn:                &grpc.ClientConn{},
		outputServiceClient: mockOutputClient,
	}

	outputServiceClient, err := c.Output()
	assert.Nil(t, err)

	stream, err := outputServiceClient.Subscribe(context.Background(), &output.Request{})
	assert.Nil(t, err)
	res, err := stream.Recv()
	assert.Nil(t, err)
	assert.Equal(t, res.Rule, "testRule", "rule does not match testRule")
	assert.Equal(t, res.Output, "testOutput", "output does not match testOutput")
}
