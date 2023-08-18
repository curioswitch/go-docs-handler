package main

import (
	"connectrpc.com/connect"
	"context"
	_ "embed"
	"fmt"
	docshandler "github.com/curioswitch/go-docs-handler"
	"github.com/curioswitch/go-docs-handler/examples/connect/greet"
	"github.com/curioswitch/go-docs-handler/examples/connect/greet/greetconnect"
	"log"
	"net/http"
	"protodescriptorset"
)

//go:embed greet/descriptors.pb
var greetDescriptors []byte

func main() {
	mux := http.NewServeMux()

	mux.Handle(greetconnect.NewGreetServiceHandler(&greetService{}))

	docsHandler, err := docshandler.New(protodescriptorset.NewPlugin(greetDescriptors))
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/docs/", http.StripPrefix("/docs", docsHandler))

	log.Fatal(http.ListenAndServe(":8080", mux))
}

type greetService struct{}

func (s *greetService) Greet(ctx context.Context, req *connect.Request[greet.GreetingRequest]) (*connect.Response[greet.GreetingResponse], error) {
	resp := &greet.GreetingResponse{
		Result: fmt.Sprintf("Hello there, %s %s", req.Msg.GetGreeting().GetFirstName(), req.Msg.GetGreeting().GetLastName()),
	}
	return connect.NewResponse(resp), nil
}
