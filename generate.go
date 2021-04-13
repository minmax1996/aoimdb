//go:generate protoc --go_out=. --go-grpc_out=. --proto_path=. --grpc-gateway_out=. --swagger_out=. --proto_path=${GOPATH}/src/ api/proto/command/command.proto
package aoimdb

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
