package tcpconnect

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"

	db "github.com/minmax1996/aoimdb/internal/aoimdb"
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
		conn:   conn,
		writer: writer,
		//outgoing:    make(chan string),
		reader:   reader,
		incoming: make(chan string),
		// stop:        make(chan bool),
		serverRef:   server,
		SessionData: SessionData{},
	}

	//go client.Write()
	go client.Read()

	return client
}

// Write writes message to the client.
func (client *Client) Write(msg interface{}) {
	var stringMsg string
	switch msg.(type) {
	case string:
		stringMsg = msg.(string)
	case []string:
		stringMsg = fmt.Sprintf("csv>%s,%s", msg.([]string)[0], msg.([]string)[1])
	case error:
		stringMsg = fmt.Sprintf("err>%s", msg.(error).Error())
	}
	logger.Info("send: " + stringMsg)
	client.writer.WriteString(stringMsg + "\n")
	client.writer.Flush()
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
			client.Write(errors.New("empty request"))
			continue
		}
		//comm:= commands.GetCommand(args[0])
		//comm.ValidateArgs(args[1:])
		//message, error := comm.CallWithArgs(args[1:])
		//client.Write( message or error
		switch cmd := args[0]; {
		case cmd == "auth":
			if len(args) == 3 {
				if err := db.DatabaseInstance.AuthentificateByUserPass(args[1], args[2]); err == nil {
					client.Write("authenticated")
					client.SessionData.authenticated = true
					continue
				}
			}
			client.Write("not authenticated")

		case cmd == "select" && client.SessionData.authenticated:
			if len(args) != 2 {
				client.Write(errors.New("not enought args to call"))
				continue
			}
			client.SessionData.selectedDatabase = args[1]
			db.DatabaseInstance.SelectDatabase(args[1])
			client.Write("selected " + args[1])
		case cmd == "get" && client.SessionData.authenticated:
			if len(args) != 2 {
				client.Write(errors.New("not enought args to call"))
				continue
			}

			selectedDatabase := client.SessionData.selectedDatabase
			key := args[1]
			if strings.Contains(args[1], ".") {
				splited := strings.Split(args[1], ".")
				selectedDatabase, key = splited[0], strings.Join(splited[1:], "")
			}

			if len(selectedDatabase) == 0 {
				client.Write(errors.New("database not selected"))
				continue
			}

			val, err := db.DatabaseInstance.Get(selectedDatabase, key)
			if err != nil {
				client.Write(err)
			} else {
				client.Write([]string{key, val.(string)})
			}
		case cmd == "set" && client.SessionData.authenticated:
			if len(args) != 3 {
				client.Write(errors.New("not enought args to call"))
				continue
			}

			selectedDatabase := client.SessionData.selectedDatabase
			key := args[1]
			value := args[2]
			if strings.Contains(args[1], ".") {
				splited := strings.Split(args[1], ".")
				selectedDatabase, key = splited[0], splited[1]
			}

			if len(selectedDatabase) == 0 {
				client.Write(errors.New("database not selected"))
				continue
			}

			err := db.DatabaseInstance.Set(selectedDatabase, key, value)
			if err != nil {
				client.Write(err)
			} else {
				client.Write("ok")
			}
		case cmd == "exit":
			client.Write("Bye")
			client.Disconnect()
			logger.Info("stop read goroutine")
			return
		case !client.SessionData.authenticated:
			client.Write("not authenticated")
		default:
			client.Write("unknown command")
		}
	}
}

// Disconnect Disconnect
func (client *Client) Disconnect() {
	client.serverRef.disconnect <- client
	client.conn.Close()
}