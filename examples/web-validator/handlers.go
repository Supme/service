package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/supme/service/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"net/http"
	"sync"
)

func PhoneSingleHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req proto.PhoneValidateRequest
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	grpcConn, err := grpc.Dial(
		server,
		grpc.WithInsecure(),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer grpcConn.Close()
	ctx := context.Background()
	validatePhone := proto.NewPhoneClient(grpcConn)
	resp, err := validatePhone.Validate(ctx, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func EmailSingleHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req proto.EmailValidateRequest
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	grpcConn, err := grpc.Dial(
		server,
		grpc.WithInsecure(),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer grpcConn.Close()
	ctx := context.Background()
	validatePhone := proto.NewEmailClient(grpcConn)
	resp, err := validatePhone.Validate(ctx, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func PhoneFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("phone_file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	fmt.Println(header.Filename)

	b := &bytes.Buffer{}
	_, err = io.Copy(b, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reader := csv.NewReader(b)
	reader.Comma = ';'
	csvData, err := reader.ReadAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	grpcConn, err := grpc.Dial(
		server,
		grpc.WithInsecure(),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer grpcConn.Close()

	ctx := context.Background()
	validatePhone := proto.NewPhoneClient(grpcConn)
	streamValidatePhone, err := validatePhone.StreamValidate(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)
	// Send
	go func() {
		for i := range csvData {
			if len(csvData[i]) < 2 {
				continue
			}
			id, number := csvData[i][0], csvData[i][1]
			req := proto.PhoneValidateRequest{Id: id, Number: number}
			//fmt.Printf("<- %+v", req)
			err := streamValidatePhone.Send(&req)
			if err != nil {
				fmt.Println("\terror happed", err)
			}
		}
		streamValidatePhone.CloseSend()
		wg.Done()
	}()
	// Receive
	go func() {
		w.Header().Set("Content-Disposition", "attachment; filename=result_"+header.Filename)
		w.Header().Set("Content-Type", "text/csv")
		writer := csv.NewWriter(w)
		writer.Comma = ';'
		for {
			r, err := streamValidatePhone.Recv()
			if err == io.EOF {
				fmt.Println("\tstream closed")
				break
			} else if err != nil {
				fmt.Println("\terror happed", err)
				break
			}
			//fmt.Printf("-> %+v", r)
			err = writer.Write([]string{r.Id, r.Canonical, fmt.Sprintf("%t", r.Valid), r.Provider, r.Error.String()})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		writer.Flush()
		wg.Done()
	}()
	wg.Wait()
}

func EmailFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("email_file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	fmt.Println(header.Filename)

	b := &bytes.Buffer{}
	_, err = io.Copy(b, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reader := csv.NewReader(b)
	reader.Comma = ';'
	csvData, err := reader.ReadAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	grpcConn, err := grpc.Dial(
		server,
		grpc.WithInsecure(),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer grpcConn.Close()

	ctx := context.Background()
	validateEmail := proto.NewEmailClient(grpcConn)
	streamValidateEmail, err := validateEmail.StreamValidate(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	// Send
	go func() {
		for i := range csvData {
			if len(csvData[i]) < 2 {
				continue
			}
			id, email := csvData[i][0], csvData[i][1]
			req := proto.EmailValidateRequest{Id: id, Email: email}
			err := streamValidateEmail.Send(&req)
			if err != nil {
				fmt.Println("\terror happed", err)
			}
		}
		streamValidateEmail.CloseSend()
		wg.Done()
	}()

	// Receive
	go func() {
		w.Header().Set("Content-Disposition", "attachment; filename=result_"+header.Filename)
		w.Header().Set("Content-Type", "text/csv")
		writer := csv.NewWriter(w)
		writer.Comma = ';'
		for {
			r, err := streamValidateEmail.Recv()
			if err == io.EOF {
				fmt.Println("\tstream closed")
				break
			} else if err != nil {
				fmt.Println("\terror happed", err)
				break
			}
			err = writer.Write([]string{r.Id, r.Canonical, fmt.Sprintf("%t", r.Valid), r.Error.String()})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		writer.Flush()
		wg.Done()
	}()
	wg.Wait()
}
