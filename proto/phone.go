//go:generate protoc -I ./ --go_out=plugins=grpc:. ./phone.proto
package proto
