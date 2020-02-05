package main

import (
	"context"
	"fmt"
	"log"

	"github.com/falcosecurity/client-go/pkg/api/version"
	"github.com/falcosecurity/client-go/pkg/client"
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

	ctx := context.Background()
	res, err := versionClient.Version(ctx, &version.Request{})
	if err != nil {
		log.Fatalf("error obtaining the Falco version: %v", err)
	}
	fmt.Printf("%v\n", res)
}
