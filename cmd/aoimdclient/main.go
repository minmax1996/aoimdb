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
	password = "admin"
	logger.Info("This is EntryPoint for client with default settings")

	conn, err := net.Dial("tcp", ":5000")
	if err != nil {
		log.Fatal(err)
	}

	connection = conn
	reader = bufio.NewReader(connection)
	writer = bufio.NewWriter(connection)
	if err := authWithCredetials(username, password); err != nil {
		log.Fatal(err)
	}

	go startListenResponses()
	commandreader := bufio.NewReader(os.Stdin)
	for {
		command, _ := commandreader.ReadString('\n')
		writer.WriteString(command)
		writer.Flush()
	}
}

func startListenResponses() {
	for {
		data, _ := reader.ReadString('\n')
		fmt.Println(data)
	}
}

func authWithCredetials(username, password string) error {
	writer.WriteString(fmt.Sprintf("auth> %s %s\n", username, password))
	writer.Flush()
	data, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	if data != "authenticated" {
		return errors.New("not authenticated")
	}
	return nil
}
