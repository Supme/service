package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/supme/service/internal/email"
	"github.com/supme/service/internal/phone"
	"github.com/supme/service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/tap"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

var (
	listenAddress string
	daemon        bool
)

func init() {
	flag.StringVar(&listenAddress, "l", ":8081", "Listen on address:port")
	flag.BoolVar(&daemon, "D", false, "Start as daemon")
	flag.Parse()

	if daemon {
		daemonize()
	}
}

func daemonize() {
	var args []string
	for i := range os.Args {
		if os.Args[i] != "-D" {
			args = append(args, os.Args[i])
		}
	}
	p := exec.Command(args[0], args[1:]...)
	p.Start()
	fmt.Println("Started pid", p.Process.Pid)
	os.Exit(0)
}

func main() {
	phoneValid, err := phone.NewValidator()
	if err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalln("can't listen port", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor),
		grpc.StreamInterceptor(authStreamInterceptor),
		grpc.InTapHandle(rateLimiter),
	)

	proto.RegisterEmailServer(server, email.NewValidator(500, 86400))
	proto.RegisterPhoneServer(server, phoneValid)

	log.Printf("starting server at %s", listenAddress)
	log.Fatal(server.Serve(l))
}

func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	reply, err := handler(ctx, req)
	fmt.Printf("time=%v token=%v\n", time.Since(start), getToken(ctx))
	return reply, err
}

func authStreamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	start := time.Now()
	err := handler(srv, stream)
	fmt.Printf("time=%v token=%v\n", time.Since(start), getToken(stream.Context()))
	return err
}

func getToken(ctx context.Context) []string {
	md, _ := metadata.FromIncomingContext(ctx)
	client := md.Get("access-token")
	return client
}

func rateLimiter(ctx context.Context, info *tap.Info) (context.Context, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	client := md.Get(":authority")
	fmt.Printf("-- rate limit check data %s authority %v\n", info.FullMethodName, client)
	return ctx, nil
}
