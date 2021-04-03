package main

import (
	"context"
	"flag"
	"net"
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/minmax1996/aoimdb/api/proto/command"
	"github.com/minmax1996/aoimdb/cmd/aoimd/grpcconnect"
	"github.com/minmax1996/aoimdb/cmd/aoimd/tcpconnect"
	"github.com/minmax1996/aoimdb/internal/aoimdb"
	"github.com/minmax1996/aoimdb/logger"
	"google.golang.org/grpc"
)

func init() {
	parseFlags()
}

func init() {
	aoimdb.InitDatabaseController()
	aoimdb.AddUser("admin", "pass")
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

func serveHttpGateway(errChan chan error) error {
	ctx := context.Background()
	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterDatabaseControllerHandlerFromEndpoint(ctx, mux, grcpPort, opts)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	go func() {
		if err := http.ListenAndServe(httpProxyPort, mux); err != nil {
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

	if err := startListenForGRPCConnects(errChan); err != nil {
		logger.Fatal(err)
		return
	}

	if err := serveHttpGateway(errChan); err != nil {
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
