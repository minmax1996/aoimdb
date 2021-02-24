package main

import (
	"context"
	"log"
	"net"

	m "github.com/minmax1996/aoimdb/api/proto/message"
	"github.com/minmax1996/aoimdb/internal/aoidb"
	"github.com/minmax1996/aoimdb/logger"
	"google.golang.org/grpc"
)

var databases map[string]*aoidb.Database

func init() {
	databases = make(map[string]*aoidb.Database)
}

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	m.UnimplementedMessangerServer
}

// SendMessage implements helloworld.GreeterServer
func (s *server) SendMessage(ctx context.Context, in *m.DefaultMessage) (*m.DefaultMessage, error) {
	logger.InfoFormat("Received: %v", in.GetQuery())
	return &m.DefaultMessage{
		Query: in.GetQuery() + ": Response",
	}, nil
}

func main() {
	logger.Info("This is EntryPoint for database service")
	databases["default"] = aoidb.NewDatabase("default")
	logger.Info(databases["default"])
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	m.RegisterMessangerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
