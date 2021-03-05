package tcpconnect

import (
	"bufio"
	"net"
	"strings"

	"github.com/minmax1996/aoimdb/logger"
)

// Client struct
type Client struct {
	conn     net.Conn
	writer   *bufio.Writer
	reader   *bufio.Reader
	incoming chan string
	outgoing chan string

	stop      chan bool
	serverRef *TCPServer

	SessionData SessionData
}

//SessionData SessionData
type SessionData struct {
	authenticated    bool
	selectedDatabase string
}

// CreateClient creates new client and starts listening
// for incoming and outgoing messages.
func CreateClient(conn net.Conn, server *TCPServer) *Client {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	client := &Client{
		conn:        conn,
		writer:      writer,
		outgoing:    make(chan string),
		reader:      reader,
		incoming:    make(chan string),
		stop:        make(chan bool),
		serverRef:   server,
		SessionData: SessionData{},
	}

	go client.Write()
	go client.Read()

	return client
}

// Write writes message to the client.
func (client *Client) Write() {
	for {
		select {
		case <-client.stop:
			logger.Info("Stop Writing gouroutine")
			client.serverRef.disconnect <- client
			break
		default:
			msg := <-client.outgoing
			logger.Info("send: " + msg)
			client.writer.WriteString(msg + "\n")
			client.writer.Flush()
		}
	}
}

// Read reads message from the client.
//TODO commands with sample checks for args, auth and etc
func (client *Client) Read() {
	for {
		msg, err := client.reader.ReadString('\n')
		if err != nil {
			client.Disconnect()
			logger.Info("stop read goroutine")
			return
		}
		args := strings.Split(strings.TrimSpace(msg), " ")
		if len(args) == 0 {
			client.outgoing <- "empty request"
			continue
		}

		switch cmd := args[0]; {
		case cmd == "auth":
			if len(args) == 3 {
				if err := client.serverRef.DatabaseController.AuthentificateByUserPass(args[1], args[2]); err == nil {
					client.outgoing <- "authenticated"
					client.SessionData.authenticated = true
					continue
				}
			}
			client.outgoing <- "not authenticated"

		case cmd == "select" && client.SessionData.authenticated:
			if len(args) != 2 {
				client.outgoing <- "not enought args to call"
				continue
			}
			client.SessionData.selectedDatabase = args[1]
			client.serverRef.DatabaseController.SelectDatabase(args[1])
			client.outgoing <- "selected " + args[1]
		case cmd == "get" && client.SessionData.authenticated:
			if len(args) != 2 {
				client.outgoing <- "not enought args to call"
				continue
			}

			if len(client.SessionData.selectedDatabase) == 0 {
				client.outgoing <- "database not selected"
				continue
			}

			val, err := client.serverRef.DatabaseController.Get(client.SessionData.selectedDatabase, args[1])
			if err != nil {
				client.outgoing <- err.Error()
			} else {
				client.outgoing <- val.(string)
			}
		case cmd == "set" && client.SessionData.authenticated:
			if len(args) != 3 {
				client.outgoing <- "not enought args to call"
				continue
			}

			if len(client.SessionData.selectedDatabase) == 0 {
				client.outgoing <- "database not selected"
				continue
			}

			err := client.serverRef.DatabaseController.Set(client.SessionData.selectedDatabase, args[1], args[2])
			if err != nil {
				client.outgoing <- err.Error()
			} else {
				client.outgoing <- "1"
			}
		case cmd == "exit":
			client.outgoing <- "Bye"
			client.Disconnect()
			logger.Info("stop read goroutine")
			return
		case !client.SessionData.authenticated:
			client.outgoing <- "not authenticated"
		default:
			client.outgoing <- "unknown command"
		}
	}
}

// Disconnect Disconnect
func (client *Client) Disconnect() {
	client.stop <- true
	client.conn.Close()
}
