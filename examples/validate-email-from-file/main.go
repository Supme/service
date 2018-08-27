package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/supme/service/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

var fileIn, fileOut, server string

func init() {
	flag.StringVar(&fileIn, "i", "./input.csv", "Input csv file contains two column: id and email")
	flag.StringVar(&fileOut, "o", "./result.csv", "Output csv file with result validate email")
	flag.StringVar(&server, "s", "127.0.0.1:8081", "Target service server")
	flag.Parse()
}

func main() {
	grpcConn, err := grpc.Dial(
		server,
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("can't connect to grpc")
	}
	defer grpcConn.Close()

	fmt.Println("Start stream email validate")
	t := time.Now()

	csvInFile, err := os.Open(fileIn)
	if err != nil {
		log.Fatal(err)
	}
	defer csvInFile.Close()
	reader := csv.NewReader(csvInFile)
	//reader.Comma = ';'
	csvData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	csvOutFile, err := os.Create(fileOut)
	if err != nil {
		log.Fatal(err)
	}
	defer csvOutFile.Close()
	writer := csv.NewWriter(csvOutFile)
	//writer.Comma = ';'
	//writer.UseCRLF = true

	ctx := context.Background()
	validateEmail := proto.NewEmailClient(grpcConn)
	streamValidateEmail, err := validateEmail.StreamValidate(ctx)
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	// Send
	go func() {
		for i := range csvData {
			id, email := csvData[i][0], csvData[i][1]
			w := proto.EmailValidateRequest{Id: id, Email: email}
			err := streamValidateEmail.Send(&w)
			if err != nil {
				fmt.Println("\terror happed", err)
			}
			//fmt.Printf("->%+v\n", w)
			// emulate network
			//time.Sleep(time.Microsecond)
		}
		streamValidateEmail.CloseSend()
		wg.Done()
	}()

	// Receive
	go func() {
		for {
			r, err := streamValidateEmail.Recv()
			if err == io.EOF {
				fmt.Println("\tstream closed")
				break
			} else if err != nil {
				fmt.Println("\terror happed", err)
				break
			}
			//fmt.Printf("<- %+v\n", r)
			err = writer.Write([]string{r.Id, r.Canonical, fmt.Sprintf("%t", r.Valid), r.Error.String()})
			if err != nil {
				log.Fatal(err)
			}
		}
		writer.Flush()
		wg.Done()
	}()
	wg.Wait()
	fmt.Printf("Check result time: %s\n", time.Since(t))

}
