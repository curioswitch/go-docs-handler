package grpcdocs

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"

	docshandler "github.com/curioswitch/go-docs-handler"
	protodocs "github.com/curioswitch/go-docs-handler/plugins/proto"
	"github.com/curioswitch/go-docs-handler/specification"
)

var errNoServices = errors.New("gRPC server does not expose any services")

func NewPlugin(s *grpc.Server, opts ...Option) docshandler.Plugin {
	config := newConfig()
	for _, opt := range opts {
		opt(config)
	}
	return &plugin{s: s, config: config}
}

type methodInfo struct {
	path string
	req  protoreflect.MessageDescriptor
	resp protoreflect.MessageDescriptor
}

type plugin struct {
	s       *grpc.Server
	config  *config
	methods []methodInfo
}

func (p *plugin) GenerateSpecification() (*specification.Specification, error) {
	var services []string
	for svc := range p.s.GetServiceInfo() {
		services = append(services, svc)
	}

	if len(services) == 0 {
		return nil, errNoServices
	}

	service, services := services[0], services[1:]
	var opts []protodocs.Option
	if p.config.serializedDescriptors != nil {
		opts = append(opts, protodocs.WithSerializedDescriptors(p.config.serializedDescriptors))
	}
	for _, svc := range services {
		opts = append(opts, protodocs.WithAdditionalService(svc))
	}
	for method, requests := range p.config.exampleRequests {
		request, requests := requests[0], requests[1:]
		opts = append(opts, protodocs.WithExampleRequests(method, request, requests...))
	}

	protoPlugin := protodocs.NewPlugin(service, opts...)

	spec, err := protoPlugin.GenerateSpecification()
	if err != nil {
		return nil, err
	}

	// A little roundabout to read from the registry again after generating the spec
	// but it keeps things relatively simple.
	for _, svc := range spec.Services {
		for _, m := range svc.Methods {
			reqDesc, _ := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(m.Parameters[0].TypeSignature.Signature()))
			respDesc, _ := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(m.ReturnTypeSignature.Signature()))
			for _, e := range m.Endpoints {
				p.methods = append(p.methods, methodInfo{
					path: e.PathMapping,
					req:  reqDesc.(protoreflect.MessageDescriptor),
					resp: respDesc.(protoreflect.MessageDescriptor),
				})
			}
		}
	}

	return spec, nil
}

type errorResponse struct {
	Code     int    `json:"code"`
	GRPCCode string `json:"grpc-code"`
	Message  string `json:"message,omitempty"`
}

func (p *plugin) AddToHandler(handler *http.ServeMux) {
	for _, m := range p.methods {
		handler.HandleFunc(m.path, func(w http.ResponseWriter, r *http.Request) {
			body := r.Body
			defer body.Close()

			reqBytes, err := io.ReadAll(body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				return
			}
			req := dynamicpb.NewMessage(m.req)
			if err := protojson.Unmarshal(reqBytes, req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			r.Header.Set("Content-Type", "application/grpc")
			reqProto, err := frameMessage(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.ProtoMajor = 2
			r.Body = io.NopCloser(bytes.NewReader(reqProto))

			rec := httptest.NewRecorder()
			p.s.ServeHTTP(rec, r)

			statusTxt := rec.Header().Get("Grpc-Status")
			code, _ := strconv.Atoi(statusTxt)
			s := status.New(codes.Code(code), rec.Header().Get("Grpc-Message"))

			w.Header().Set("Content-Type", "application/json")

			if s.Code() != codes.OK {
				resp := errorResponse{
					Code:     code,
					GRPCCode: s.Code().String(),
					Message:  s.Message(),
				}
				respBytes, err := json.Marshal(resp)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(grpcCodeToHTTPStatus(s.Code()))
				_, _ = w.Write(respBytes)
				return
			}

			framed := rec.Body.Bytes()
			// Go ahead and read out the size even though we don't expect more than one message.
			sz := binary.BigEndian.Uint32(framed[1:5])
			respProto := framed[5 : 5+sz]
			resp := dynamicpb.NewMessage(m.resp)
			if err := proto.Unmarshal(respProto, resp); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			respBytes, err := protojson.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(respBytes)
		})
	}
}

func frameMessage(msg proto.Message) ([]byte, error) {
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("marshal message to binary: %w", err)
	}

	framed := make([]byte, len(msgBytes)+5)
	framed[0] = 0
	binary.BigEndian.PutUint32(framed[1:5], uint32(len(msgBytes)))
	copy(framed[5:], msgBytes)
	return framed, nil
}

func grpcCodeToHTTPStatus(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return 499 // CLIENT_CLOSED_REQUEST
	case codes.Unknown:
	case codes.Internal:
	case codes.DataLoss:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
	case codes.FailedPrecondition:
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
	case codes.Aborted:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	}
	return http.StatusInternalServerError
}
