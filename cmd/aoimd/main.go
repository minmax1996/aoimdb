package main

import (
	"flag"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/minmax1996/aoimdb/cmd/aoimd/tcpconnect"
	"github.com/minmax1996/aoimdb/internal/aoimdb/database"
	"github.com/minmax1996/aoimdb/pkg/logger"
)

func init() {
	parseFlags()
	database.InitDatabaseController()
	_ = database.AddUser("admin", "pass")
}

const (
	tcpPort       = ":1593" // last digit of "aoim": a=1 o=15 i=9 m=13
	grcpPort      = ":50051"
	httpProxyPort = ":8081"
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
	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir(assetsPath)))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(assetsPath+"/static/"))))
	go func() {
		if err := http.ListenAndServe(":3000", r); err != nil {
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

	if err := startWebUI(errChan); err != nil {
		logger.Fatal(err)
		return
	}

	//TODO MAnual Backup on stop app

	//main loop for errors
	for err := range errChan {
		logger.Error(err)
	}
}
