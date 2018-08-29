package main

import (
	"encoding/json"
	"github.com/supme/service/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net/http"
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

}

func EmailFileHandler(w http.ResponseWriter, r *http.Request) {

}
