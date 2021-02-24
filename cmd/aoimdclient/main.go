package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/minmax1996/aoimdb/logger"
	"golang.org/x/net/http2"
)

func main() {
	logger.Info("This is EntryPoint for client with default settings")
	client := &http.Client{}

	// Use the proper transport in the client
	client.Transport = &http2.Transport{}

	// Perform the request
	resp, err := client.Post("https://localhost:9191/hello/sayHello", "text/plain", bytes.NewBufferString("Hello Go!"))
	if err != nil {
		log.Fatalf("Failed get: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed reading response body: %s", err)
	}
	fmt.Printf("Got response %d: %s %s", resp.StatusCode, resp.Proto, string(body))
}
