package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/supme/service/proto"
	"google.golang.org/grpc"
	"log"
)

var server, phonenumber string

func init() {
	flag.StringVar(&server, "s", "127.0.0.1:8081", "Target service server")
	flag.StringVar(&phonenumber, "p", "+71234567890", "Phone for validate")
	flag.Parse()
}

func main() {
	grpcConn, err := grpc.Dial(
		server,
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal("can't connect to grpc")
	}
	defer grpcConn.Close()

	ctx := context.Background()
	validatePhone := proto.NewPhoneClient(grpcConn)

	resp, err := validatePhone.Validate(ctx, &proto.PhoneValidateRequest{Id: "test", Number: phonenumber})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Id: %s Canonical: %s Valid: %t Provider: %s Error: %s\n", resp.Id, resp.Canonical, resp.Valid, resp.Provider, resp.Error)
}
