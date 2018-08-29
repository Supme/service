package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var (
	listen, server string
	daemon         bool
)

func init() {
	flag.StringVar(&listen, "l", ":8080", "Listen interface:port")
	flag.StringVar(&server, "s", "127.0.0.1:8081", "Target service server")
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
	http.Handle("/", http.FileServer(http.Dir("./assets")))
	http.HandleFunc("/phone/single", PhoneSingleHandler)
	http.HandleFunc("/email/single", EmailSingleHandler)
	http.HandleFunc("/phone/file", PhoneFileHandler)
	http.HandleFunc("/email/file", EmailFileHandler)
	fmt.Printf("Listen on %s\n", listen)
	log.Fatal(http.ListenAndServe(listen, nil))
}
