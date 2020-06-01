package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/falcosecurity/client-go/pkg/api/outputs"
	"github.com/falcosecurity/client-go/pkg/client"
	"github.com/gogo/protobuf/jsonpb"
)

func main() {
	//Set up a connection to the server.
	c, err := client.NewForConfig(context.Background(), &client.Config{
		UnixSocketPath: "unix:///var/run/falco.sock",
	})
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer c.Close()
	outputsClient, err := c.Outputs()
	if err != nil {
		log.Fatalf("unable to obtain an output client: %v", err)
	}

	ctx := context.Background()
	fcs, err := outputsClient.Sub(ctx)
	if err != nil {
		log.Fatalf("could not subscribe: %v", err)
	}

	go func() {
		for {
			fcs.Send(&outputs.Request{})
			time.Sleep(time.Second * 5)
		}
	}()
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
