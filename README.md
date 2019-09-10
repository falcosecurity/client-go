# Falco Go Client

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
    "google.golang.org/grpc"
)

func main() {
    c, err := client.NewForConfig(&client.Config{
        Target: "localhost:5060",
        Options: []grpc.DialOption{
            grpc.WithInsecure(),
        },
    })
}
```


### Falco outputs subscribe

```go
outputClient, err := c.Output()
if err != nil {
    log.Fatalf("unable to obtain an output client: %v", err)
}

ctx := context.Background()
// Keepalive true means that the client will wait indefinitely for new events to come
// Use keepalive false if you only want to receive the accumulated events and stop
fcs, err := outputClient.Subscribe(ctx, &output.FalcoOutputRequest{Keepalive: true})
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

## Full Examples

- [Output events example](examples/output/main.go)
