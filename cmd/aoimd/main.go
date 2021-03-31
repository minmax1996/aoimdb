package main

import (
	"flag"
	"net"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/minmax1996/aoimdb/cmd/aoimd/tcpconnect"
	"github.com/minmax1996/aoimdb/internal/aoimdb"
	"github.com/minmax1996/aoimdb/logger"
)

func init() {
	parseFlags()
	aoimdb.InitDatabaseController()
	aoimdb.AddUser("admin", "pass")
}

const (
	tcpPort = ":1593" // last digit of "aoim": a=1 o=15 i=9 m=13
)

var (
	assetsPath = "./ui/build"
)

func parseFlags() {
	flag.StringVar(&assetsPath, "assets_path", "./ui/build", "a string var for username")
	flag.Parse()
}

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

func startWebUI(errChan chan error) error {
	router := gin.Default()
	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile(assetsPath, true)))

	// Start and run the server
	go func(errorChannel chan error) {
		err := router.Run(":3000")
		if err != nil {
			errorChannel <- err
			return
		}
	}(errChan)

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

	if err := startWebUI(errChan); err != nil {
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
