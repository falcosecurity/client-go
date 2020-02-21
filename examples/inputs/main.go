package main

import (
	"context"
	"fmt"
	"log"

	"github.com/falcosecurity/client-go/pkg/api/inputs"
	"github.com/falcosecurity/client-go/pkg/client"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	// Set up a connection to the server.
	c, err := client.NewForConfig(&client.Config{
		Hostname:   "localhost",
		Port:       5060,
		CertFile:   "/tmp/client.crt",
		KeyFile:    "/tmp/client.key",
		CARootFile: "/tmp/ca.crt",
	})
	if err != nil {
		log.Fatalf("unable to create a Falco client: %v", err)
	}
	defer c.Close()

	inputsClient, err := c.Inputs()
	if err != nil {
		log.Fatalf("unable to obtain a inputs client: %v", err)
	}

	event0 := &any.Any{Value: []byte(`{"num":0"}`)}
	event1 := &any.Any{Value: []byte(`{"num":1"}`)}
	pack0, err := ptypes.MarshalAny(event0)
	if err != nil {
		log.Fatalf("unable to pack event 0")
	}
	pack1, err := ptypes.MarshalAny(event1)
	if err != nil {
		log.Fatalf("unable to pack event 1")
	}

	req := &inputs.Request{
		Events: []*any.Any{pack0, pack1},
	}

	var header, trailer metadata.MD

	res, err := inputsClient.Input(
		context.Background(),
		req,
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	if err != nil {
		log.Fatalf("error obtaining a reply: %v", err)
	}
	out, err := (&jsonpb.Marshaler{}).MarshalToString(res)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
	fmt.Println(header)
	fmt.Println(trailer)
}
