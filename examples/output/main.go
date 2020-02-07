package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/falcosecurity/client-go/pkg/api/output"
	"github.com/falcosecurity/client-go/pkg/client"
	"github.com/gogo/protobuf/jsonpb"
)

func main() {
	//Set up a connection to the server.
	c, err := client.NewForConfig(&client.Config{
		Hostname:   "localhost",
		Port:       5060,
		CertFile:   "/tmp/client.crt",
		KeyFile:    "/tmp/client.key",
		CARootFile: "/tmp/ca.crt",
	})
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer c.Close()
	outputClient, err := c.Output()
	if err != nil {
		log.Fatalf("unable to obtain an output client: %v", err)
	}

	ctx := context.Background()
	// Keepalive true means that the client will wait indefinitely for new events to come
	// Use keepalive false if you only want to receive the accumulated events and stop
	fcs, err := outputClient.Subscribe(ctx, &output.Request{Keepalive: true})
	if err != nil {
		log.Fatalf("could not subscribe: %v", err)
	}

	for {
		res, err := fcs.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error closing stream after EOF: %v", err)
		}
		out, err := (&jsonpb.Marshaler{}).MarshalToString(res)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(out)
	}
}
