package main

import (
	"flag"
	"fmt"
	"github.com/supme/service/internal/email"
	"github.com/supme/service/internal/translit"
	"github.com/supme/service/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/exec"
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
	l, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalln("can't listen port", err)
	}

	server := grpc.NewServer()

	proto.RegisterTransliterationServer(server, translit.NewTr())
	proto.RegisterEmailServer(server, email.NewValidator(50, 30))

	log.Printf("starting server at %s", listenAddress)
	log.Fatal(server.Serve(l))
}
