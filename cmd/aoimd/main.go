package main

import (
	"net"

	"github.com/minmax1996/aoimdb/internal/aoidb"
	"github.com/minmax1996/aoimdb/logger"
)

var (
	databases map[string]*aoidb.Database
	server    *Chat
	errChan   chan error
)

func init() {
	databases = make(map[string]*aoidb.Database)
	server = CreateChat()
	errChan = make(chan error)
}

const (
	port = ":50051"
)

func startListenForUnprotectedTCP(errChan chan error) error {
	logger.Info("Starting listening for connections...")

	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		logger.Error(err)
		return err
	}

	go func(ec chan error) {
		for { // listen for connections
			conn, err := listener.Accept()
			if err != nil {
				logger.Error(err)
				ec <- err
				continue
			}
			logger.Info("Connected")
			server.Connect(conn)
		}
	}(errChan)

	return nil
}

func main() {
	logger.Info("This is EntryPoint for database service")
	databases["default"] = aoidb.NewDatabase("default")
	logger.Info(databases["default"])

	if err := startListenForUnprotectedTCP(errChan); err != nil {
		logger.Fatal(err)
		return
	}

	for {
		select {
		case err := <-errChan:
			logger.Error(err)
		}
	}
}
