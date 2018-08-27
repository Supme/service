//go:generate protoc -I ./ --go_out=plugins=grpc:. ./email.proto
package proto
