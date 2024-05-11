package main

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	ins "grpcgateway/insecure"
	server "grpcgateway/server"

	"github.com/felixge/httpsnoop"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"

	gw "grpcgateway/proto"
)

var (
	HTTPS      = false
	GRPC       = "0.0.0.0:8080"
	DNS        = "dns:///0.0.0.0:8080"
	GATEWAY    = "0.0.0.0:8090"
	GATEWAYMUX *runtime.ServeMux
)

func withLogger(handler http.Handler) http.Handler {
	// the create a handler
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// pass the handler to httpsnoop to get http status and latency
		m := httpsnoop.CaptureMetrics(handler, writer, request)
		// printing exracted data
		log.Printf("http[%d]-- %s -- %s\n", m.Code, m.Duration, request.URL.Path)
	})
}

var allowedHeaders = map[string]struct{}{
	"x-request-id": {},
}

func isHeaderAllowed(s string) (string, bool) {
	// check if allowedHeaders contain the header
	if _, isAllowed := allowedHeaders[s]; isAllowed {
		// send uppercase header
		return strings.ToUpper(s), true
	}
	// if not in the allowed header, don't send the header
	return s, false
}

func main() {
	grlog := grpclog.NewLoggerV2(os.Stdout, io.Discard, io.Discard)
	grpclog.SetLoggerV2(grlog)

	lis := listener()
	registerServer(lis)
	conn := dial()
	mux()
	registerHandler(conn)
	serve()
}

func listener() net.Listener {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", GRPC)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	return lis
}

func registerServer(lis net.Listener) {
	var GRPC_SERVER *grpc.Server

	if HTTPS {
		GRPC_SERVER = grpc.NewServer(
			grpc.Creds(credentials.NewServerTLSFromCert(&ins.Cert)),
		)
	} else {
		GRPC_SERVER = grpc.NewServer()
	}

	// Attach the Greeter service to the server
	newServer := server.New()
	gw.RegisterGreeterServer(GRPC_SERVER, newServer)
	gw.RegisterUserServiceServer(GRPC_SERVER, newServer)
	// Serve gRPC server
	log.Printf("Serving gRPC on %s", GRPC)
	go func() {
		log.Fatalln(GRPC_SERVER.Serve(lis))
	}()
}

func dial() *grpc.ClientConn {
	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests

	var cred credentials.TransportCredentials

	if HTTPS {
		cred = credentials.NewClientTLSFromCert(ins.CertPool, "")
	} else {
		cred = insecure.NewCredentials()
	}

	conn, err := grpc.DialContext(
		context.Background(),
		DNS,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(cred),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	return conn
}

func mux() {
	GATEWAYMUX = runtime.NewServeMux(
		runtime.WithOutgoingHeaderMatcher(isHeaderAllowed),
		runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
			header := r.Header.Get("Authorization")
			// send all the headers received from the client
			md := metadata.Pairs("auth", header)
			return md
		}),
		runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
			// creating a new HTTTPStatusError with a custom status, and passing error
			newError := runtime.HTTPStatusError{
				HTTPStatus: 400,
				Err:        err,
			}
			// using default handler to do the rest of heavy lifting of marshaling error and adding headers
			runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, writer, request, &newError)
		}),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),
	)
}

func registerHandler(conn *grpc.ClientConn) {
	err := gw.RegisterGreeterHandler(context.Background(), GATEWAYMUX, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	err = gw.RegisterUserServiceHandler(context.Background(), GATEWAYMUX, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
}

func serve() {
	gwServer := &http.Server{
		Addr:    GATEWAY,
		Handler: withLogger(GATEWAYMUX),
	}

	if HTTPS {
		gwServer.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{ins.Cert},
		}
		log.Printf("Serving gRPC-Gateway on https://%s", GATEWAY)
		log.Fatalln(gwServer.ListenAndServeTLS("", ""))
	} else {
		log.Printf("Serving gRPC-Gateway on http://%s", GATEWAY)
		log.Fatalln(gwServer.ListenAndServe())
	}
}
