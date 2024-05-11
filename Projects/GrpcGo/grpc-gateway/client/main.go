package main

import (
	"log"

	gw "grpcgateway/proto"

	ins "grpcgateway/insecure"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	HTTPS   = false
	GRPC    = "0.0.0.0:8080"
	DNS     = "dns:///0.0.0.0:8080"
	GATEWAY = "0.0.0.0:8090"
)

// Authentication holds the login/password
type Authentication struct {
	Login    string
	Password string
}

// GetRequestMetadata gets the current request metadata
func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	log.Printf("Login : %s \nPassword : %s\n", a.Login, a.Password)
	return map[string]string{
		"login":    a.Login,
		"password": a.Password,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security
func (a *Authentication) RequireTransportSecurity() bool {
	return true
}

func main() {
	var cred credentials.TransportCredentials

	if HTTPS {
		cred = credentials.NewClientTLSFromCert(ins.CertPool, "")
	} else {
		cred = insecure.NewCredentials()
	}

	// Setup the login/pass
	auth := Authentication{
		Login:    "Grpc",
		Password: "Gateway",
	}

	var conn *grpc.ClientConn
	var err error

	if HTTPS {
		// Initiate a connection with the server
		conn, err = grpc.Dial(
			"0.0.0.0:8080",
			grpc.WithTransportCredentials(cred),
			grpc.WithPerRPCCredentials(&auth),
		)
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
	} else {
		// Initiate a connection with the server
		conn, err = grpc.Dial(
			"0.0.0.0:8080",
			grpc.WithTransportCredentials(cred),
		)
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
	}

	defer conn.Close()

	c := gw.NewGreeterClient(conn)

	response, err := c.SayHello(context.Background(), &gw.HelloRequest{Name: "GrpcGateway"})
	if err != nil {
		log.Fatalf("error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Message)
}
