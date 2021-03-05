package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/minmax1996/aoimdb/logger"
)

var (
	username   string
	password   string
	connection net.Conn
	reader     *bufio.Reader
	writer     *bufio.Writer
)

func main() {
	username = "admin"
	password = "pass"
	logger.Info("This is EntryPoint for client with default settings")

	conn, err := net.Dial("tcp", ":1593")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	connection = conn
	reader = bufio.NewReader(connection)
	writer = bufio.NewWriter(connection)
	if err := authWithCredetials(username, password); err != nil {
		log.Fatal(err)
	}

	go startListenResponses()
	commandreader := bufio.NewReader(os.Stdin)
	for {
		command, err := commandreader.ReadString('\n')
		if err != nil {
			os.Exit(1)
			break
		}
		writer.WriteString(command)
		writer.Flush()
	}
	os.Exit(0)
}

func startListenResponses() {
	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			os.Exit(1)
		}
		fmt.Println("<" + data)
	}
}

func authWithCredetials(username, password string) error {
	writer.WriteString(fmt.Sprintf("auth %s %s\n", username, password))
	writer.Flush()
	data, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	logger.Info("auth response: " + data)
	if data != "authenticated\n" {
		return errors.New("not authenticated")
	}
	return nil
}
