package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"

	docshandler "github.com/curioswitch/go-docs-handler"
	"github.com/curioswitch/go-docs-handler/examples/grpc/greet"
	grpcdocs "github.com/curioswitch/go-docs-handler/plugins/grpc"
)

//go:embed greet/descriptors.pb
var greetDescriptors []byte

func main() {
	server := grpc.NewServer()
	defer server.Stop()

	greet.RegisterGreetServiceServer(server, service{})

	docsHandler, err := docshandler.New(grpcdocs.NewPlugin(server,
		grpcdocs.WithSerializedDescriptors(greetDescriptors),
		grpcdocs.WithExampleRequests("greet.GreetService/Greet",
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
			})))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := http.ListenAndServe(":8081", docsHandler); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	list, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Serve(list); err != nil {
		log.Fatal(err)
	}
}

type service struct {
	greet.UnimplementedGreetServiceServer
}

func (service) Greet(ctx context.Context, req *greet.GreetingRequest) (*greet.GreetingResponse, error) {
	res := "Who are you?"
	g := req.Greeting
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

	return &greet.GreetingResponse{
		Result: res,
	}, nil
}
