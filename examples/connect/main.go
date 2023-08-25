package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	docshandler "github.com/curioswitch/go-docs-handler"
	"github.com/curioswitch/go-docs-handler/examples/connect/greet"
	"github.com/curioswitch/go-docs-handler/examples/connect/greet/greetconnect"
	"github.com/curioswitch/go-docs-handler/plugins/protodescriptorset"
)

//go:embed greet/descriptors.pb
var greetDescriptors []byte

func main() {
	mux := http.NewServeMux()

	mux.Handle(greetconnect.NewGreetServiceHandler(&greetService{}))

	docsHandler, err := docshandler.New(protodescriptorset.NewPlugin(greetDescriptors,
		protodescriptorset.WithExampleRequests(greetconnect.GreetServiceGreetProcedure[1:], []proto.Message{
			&greet.GreetingRequest{
				Greeting: &greet.Greeting{
					Name: &greet.Greeting_Nickname{
						Nickname: "Choko",
					},
				},
			},
			&greet.GreetingRequest{
				Greeting: &greet.Greeting{
					Name: &greet.Greeting_FullName_{
						FullName: &greet.Greeting_FullName{
							FirstName: "Choko",
							LastName:  "Switch",
						},
					},
				},
			},
			&greet.GreetingRequest{
				Greeting: &greet.Greeting{
					Name: &greet.Greeting_KnownName_{
						KnownName: greet.Greeting_BOB,
					},
				},
			},
		})))
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/docs/", http.StripPrefix("/docs", docsHandler))

	log.Fatal(http.ListenAndServe(":8080", mux))
}

type greetService struct{}

func (s *greetService) Greet(ctx context.Context, req *connect.Request[greet.GreetingRequest]) (*connect.Response[greet.GreetingResponse], error) {
	g := req.Msg.GetGreeting()

	res := "Who are you?"
	switch n := g.GetName().(type) {
	case *greet.Greeting_Nickname:
		res = fmt.Sprintf("Hello there, %s", n.Nickname)
	case *greet.Greeting_FullName_:
		res = fmt.Sprintf("Hello there, %s %s", n.FullName.GetFirstName(), n.FullName.GetLastName())
	case *greet.Greeting_KnownName_:
		switch n.KnownName {
		case greet.Greeting_BOB:
			res = "Hello there, Bob"
		case greet.Greeting_ALICE:
			res = "Hello there, Alice"
		}
	}

	resp := &greet.GreetingResponse{
		Result: res,
	}
	return connect.NewResponse(resp), nil
}
