package main

import (
	"flag"
	"fmt"
	"github.com/supme/service/internal/email"
	"github.com/supme/service/internal/phone"
	"github.com/supme/service/internal/translit"
	"github.com/supme/service/proto"
	"google.golang.org/grpc"
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

	phones := []string{
		"+79066460826",
		"(910)563-82-56",
		"8-900-905-57-69",
		"123",
		"4912984710",
		"8002242222",
		"+8940984093",
		"4952242222",
		"81234567890",
		"88005555550",
		"8000000000",
		"+19004561244",
	}
	start := time.Now()
	for i := range phones {
		canonical, provider, err := phoneValid.Check(phones[i])
		if err != nil {
			fmt.Printf("Phone %s check error: '%s'\n", phones[i], err)
			continue
		}
		fmt.Printf("Phone %s\tcanonical format %s\tprovider %s\n", phones[i], canonical, provider)
	}
	fmt.Printf("Time to find %s\n", time.Since(start))

	return

	l, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalln("can't listen port", err)
	}

	server := grpc.NewServer()

	proto.RegisterTransliterationServer(server, translit.NewTr())
	proto.RegisterEmailServer(server, email.NewValidator(100, 3600))
	proto.RegisterPhoneServer(server, phoneValid)

	log.Printf("starting server at %s", listenAddress)
	log.Fatal(server.Serve(l))
}
