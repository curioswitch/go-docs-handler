# go-docs-handler

go-docs-handler is a `net/http` handler that serves documentation for API services written in Go.
The client is the same as the excellent [Armeria][0] docs client, with the server-side handling of
specification requests reimplemented with `net/http`.

## Setup

`go-docs-handler` aims to support different frameworks generically by using framework-specific
plugins. Currently only gRPC using either [Connect][1] or [official gRPC][2] are supported.

### Connect

Support for Connect is provided by the [proto](./plugins/proto) plugin, which reads definitions from the global
protobuf registry. It does not actually have any dependencies on Connect itself and can be
used with other implementations of protobuf-based RPC. If anyone has a framework that works with it,
let us know so we can add docs and examples for it.

```bash
go get github.com/curioswitch/go-docs-handler/plugins/proto
```

```go
import "github.com/curioswitch/go-docs-handler/plugins/proto"

func main() {
	mux := http.NewServeMux()
	
	// Register the connect handler as usual
	mux.Handle(greetconnect.NewGreetServiceHandler(&greetService{}))
	
	docsHandler, err := docshandler.New(protodocs.NewPlugin(greetconnect.GreetServiceName))
	if err != nil {
		panic(err)
	}
	
	// Register the docs handler onto the same mux. It's recommended to give it a prefix.
	mux.Handle("/docs/", http.StripPrefix("/docs", docsHandler))
	
	http.ListenAndServe(":8080", mux)
}
```

Documentation will be served on `http://localhost:8080/docs/` and debug requests will be served
by the registered handler. Note that services are registered to documentation by name as there
is no way to introspect what connect services have been registered on a mux. If the service name(s)
used when initializing the plugin do not match the actual handlers, documentation will be served
for the unrelated service and debug requests will fail.

Also see a full example [here](./examples/connect).

### gRPC

Support for gRPC is provided by the [gRPC](./plugins/grpc) plugin, which reads registered services
from a `grpc.Server` and uses the global registry to resolve their details. The gRPC server uses a
custom HTTP stack and protocol and cannot be used to serve docs or debug requests. The docs handler
should be registered on a separate HTTP port, which will also serve debug requests from the browser,
proxying them to the gRPC server. Currently, the docs handler must be served at the root when using
the gRPC plugin.

```bash
go get github.com/curioswitch/go-docs-handler/plugins/grpc
```

```go
import "github.com/curioswitch/go-docs-handler/plugins/grpc"

func main() {
	server := grpc.NewServer()
	greet.RegisterGreetServiceServer(server, service{})
	
	mux := http.NewServeMux()
	
	
	docsHandler, err := docshandler.New(protodocs.NewPlugin(server))
	if err != nil {
		panic(err)
	}
	
	go func() {
		http.ListenAndServe(":8081", docsHandler)
	}

	list, _ := net.Listen("tcp", ":8080")
	server.Serve(list)
}
```

Documentation will be served on `http://localhost:8081/` and debug requests will be served
by the registered handler, proxied through the docs handler.

Also see a full example [here](./examples/grpc).

[0]: https://armeria.dev
[1]: https://connectrpc.com/docs/go
[2]: https://github.com/grpc/grpc-go
