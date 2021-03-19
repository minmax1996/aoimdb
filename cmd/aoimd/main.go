package main

import (
	"net"

	pb "github.com/minmax1996/aoimdb/api/proto/command"
	"github.com/minmax1996/aoimdb/cmd/aoimd/grpcconnect"
	"github.com/minmax1996/aoimdb/cmd/aoimd/tcpconnect"
	"github.com/minmax1996/aoimdb/internal/aoimdb"
	"github.com/minmax1996/aoimdb/logger"
	"google.golang.org/grpc"
)

func init() {
	aoimdb.InitDatabaseController()
	aoimdb.AddUser("admin", "pass")
}

const (
	tcpPort  = ":1593" // last digit of "aoim": a=1 o=15 i=9 m=13
	grcpPort = ":50051"
)

func startListenForTCPConnects(errChan chan error) error {
	server := tcpconnect.CreateTCPServer(errChan)
	listener, err := net.Listen("tcp", tcpPort)
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.InfoFormat("Start listening tcp connections on: %s", tcpPort)

	go func(errorChannel chan error) {
		for { // listen for connections
			conn, err := listener.Accept()
			if err != nil {
				errorChannel <- err
				continue
			}
			logger.Info("Connected")
			server.Connect(conn)
		}
	}(errChan)

	return nil
}

func startListenForGRPCConnects(errChan chan error) error {
	listener, err := net.Listen("tcp", grcpPort)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterDatabaseControllerServer(s, &grpcconnect.Server{})

	go func() {
		if err := s.Serve(listener); err != nil {
			errChan <- err
		}
	}()
	return nil
}

func main() {
	// main error chan
	errChan := make(chan error)
	logger.Info("This is EntryPoint for database service")

	if err := startListenForTCPConnects(errChan); err != nil {
		logger.Fatal(err)
		return
	}

	if err := startListenForGRPCConnects(errChan); err != nil {
		logger.Fatal(err)
		return
	}

	//TODO MAnual Backup on stop app

	//main loop for errors
	for {
		select {
		case err := <-errChan:
			logger.Error(err)
		}
	}
}
