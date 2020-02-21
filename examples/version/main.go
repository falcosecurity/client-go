package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/falcosecurity/client-go/pkg/api/version"
	"github.com/falcosecurity/client-go/pkg/client"
	"github.com/gogo/protobuf/jsonpb"
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
	versionClient, err := c.Version()
	if err != nil {
		log.Fatalf("unable to obtain a version client: %v", err)
	}

	var header, trailer metadata.MD
	res, err := versionClient.Version(
		context.Background(),
		&version.Request{},
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	if err != nil {
		log.Fatalf("error obtaining the Falco version: %v", err)
	}
	out, err := (&jsonpb.Marshaler{}).MarshalToString(res)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
	// Header metadata
	headerString, err := json.Marshal(header)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(headerString))
	// Trailer metadata
	trailerString, err := json.Marshal(trailer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(trailerString))
}
