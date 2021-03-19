//go:generate protoc --go_out=. --go-grpc_out=. --proto_path=. --proto_path=${GOPATH}/src/ api/proto/command/command.proto
package aoimdb
