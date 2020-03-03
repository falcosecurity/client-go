# Falco Go Client
[![GoDoc](https://godoc.org/github.com/falcosecurity/client-go/pkg/client?status.svg)](https://godoc.org/github.com/falcosecurity/client-go/pkg/client)

> Go client and SDK for Falco

## Install

```bash
go get -u github.com/falcosecurity/client-go
```

## Usage

### Client creation

```go
package main

imports(
    "github.com/falcosecurity/client-go/pkg/client"
)

func main() {
    c, err := client.NewForConfig(&client.Config{
        Hostname:   "localhost",
        Port:       5060,
        CertFile:   "/etc/falco/certs/client.crt",
        KeyFile:    "/etc/falco/certs/client.key",
        CARootFile: "/etc/falco/certs/ca.crt",
    })
}
```

### Falco outputs API

```go
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
    fmt.Printf("rule: %s\n", res.Rule)
}
```

### Falco version API

```go
// Set up a connection to the server.
c, err := client.NewForConfig(&client.Config{
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
fmt.Printf("%v\n", res)
```

## Full Examples

- [Output events example](examples/output/main.go)
- [Version example](examples/version/main.go)

## Update protos

Perform the following edits to the Makefile:

1. Update the `PROTOS` array with the destination path of the `.proto` file.
2. Update the `PROTO_URLS` array with the URL from which to download it.
3. Update thr `PROTO_SHAS` array with the SHA256 sum of the file to download.
4. Execute the following commands:

```console
make clean
make protos
```
