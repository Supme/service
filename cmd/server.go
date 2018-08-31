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

	//// ------------- test phone validator -----------
	//phones := []string{
	//	// Kazahstan
	//	"+77122123456",
	//	"+73362212345",
	//}
	//start := time.Now()
	//for i := range phones {
	//	canonical, provider, err := phoneValid.Check(phones[i])
	//	if err != nil {
	//		fmt.Printf("Phone %s check error: '%s'\n", phones[i], err)
	//		continue
	//	}
	//	fmt.Printf("Phone %s\tcanonical format %s\tprovider %s\n", phones[i], canonical, provider)
	//}
	//fmt.Printf("Time to find %s\n", time.Since(start))
	//// ----------------------------------------------

	l, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalln("can't listen port", err)
	}

	server := grpc.NewServer()

	proto.RegisterTransliterationServer(server, translit.NewTr())
	proto.RegisterEmailServer(server, email.NewValidator(500, 86400))
	proto.RegisterPhoneServer(server, phoneValid)

	log.Printf("starting server at %s", listenAddress)
	log.Fatal(server.Serve(l))
}
