//go:generate protoc -I ./ --go_out=plugins=grpc:. ./translit.proto
package proto
