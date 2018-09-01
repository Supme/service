package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/supme/service/proto"
	"google.golang.org/grpc"
	"log"
	"time"
)

var server, phoneNumber string

func init() {
	flag.StringVar(&server, "s", "127.0.0.1:8081", "Target service server")
	flag.StringVar(&phoneNumber, "p", "+71234567890", "Phone for validate")
	flag.Parse()
}

func main() {
	grpcConn, err := grpc.Dial(
		server,
		grpc.WithUnaryInterceptor(timingInterceptor),
		grpc.WithPerRPCCredentials(&tokenAuth{"MySecureToken"}),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal("can't connect to grpc")
	}
	defer grpcConn.Close()

	ctx := context.Background()
	validatePhone := proto.NewPhoneClient(grpcConn)

	resp, err := validatePhone.Validate(ctx, &proto.PhoneValidateRequest{Id: "test", Number: phoneNumber})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Id: %s Canonical: %s Valid: %t Provider: %s Error: %s\n", resp.Id, resp.Canonical, resp.Valid, resp.Provider, resp.Error)
}

type tokenAuth struct {
	Token string
}

func (t *tokenAuth) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"access-token": t.Token}, nil
}

func (t *tokenAuth) RequireTransportSecurity() bool {
	return false
}

func timingInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	fmt.Printf("time=%v\n", time.Since(start))
	return err

}
