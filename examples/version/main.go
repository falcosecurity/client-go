package main

import (
	"context"
	"fmt"
	"log"

	"github.com/falcosecurity/client-go/pkg/api/version"
	"github.com/falcosecurity/client-go/pkg/client"
	"github.com/gogo/protobuf/jsonpb"
)

func main() {
	// Set up a connection to the server.
	c, err := client.NewForConfig(context.Background(), &client.Config{
		Hostname:   "localhost",
		Port:       5060,
		CertFile:   "/etc/falco/certs/client.crt",
		KeyFile:    "/etc/falco/certs/client.key",
		CARootFile: "/etc/falco/certs/ca.crt",
	})
	if err != nil {
		log.Fatalf("unable to create a Falco client: %v", err)
	}
	defer c.Close()
	versionClient, err := c.Version()
	if err != nil {
		log.Fatalf("unable to obtain a version client: %v", err)
	}

	ctx := context.Background()
	res, err := versionClient.Version(ctx, &version.Request{})
	if err != nil {
		log.Fatalf("error obtaining the Falco version: %v", err)
	}
	out, err := (&jsonpb.Marshaler{}).MarshalToString(res)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
}
