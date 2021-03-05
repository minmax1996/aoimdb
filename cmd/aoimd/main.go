package main

import (
	"net"

	"github.com/minmax1996/aoimdb/cmd/aoimd/tcpconnect"
	"github.com/minmax1996/aoimdb/internal/aoimdb"
	"github.com/minmax1996/aoimdb/logger"
)

var (
	databases *aoimdb.DatabaseController
	errChan   chan error
)

func init() {
	errChan = make(chan error)
	databases = aoimdb.NewDatabaseController()
	databases.AddUser("admin", "pass")
}

const (
	tcpPort = ":1593" // last digit of "aoim": a=1 o=15 i=9 m=13
)

func startListenForUnprotectedTCP(errChan chan error) error {
	server := tcpconnect.CreateTCPServer(databases)
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

func main() {
	logger.Info("This is EntryPoint for database service")
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
