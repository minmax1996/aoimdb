package tcpconnect

import (
	"bufio"
	"net"
	"strings"

	"github.com/minmax1996/aoimdb/api/commands"
	"github.com/minmax1996/aoimdb/api/msg_protocol"
	"github.com/minmax1996/aoimdb/logger"
	"github.com/vmihailenco/msgpack/v5"
)

// Client struct
type Client struct {
	conn     net.Conn
	writer   *bufio.Writer
	reader   *bufio.Reader
	incoming chan string

	disconnect chan *Client

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
	client := &Client{
		conn:        conn,
		writer:      bufio.NewWriter(conn),
		reader:      bufio.NewReader(conn),
		incoming:    make(chan string),
		disconnect:  server.disconnect,
		SessionData: SessionData{},
	}

	go client.Read()

	return client
}

// WriteError wrapper around Write to write only error
func (client *Client) WriteError(err error) {
	client.Write(&msg_protocol.MsgPackRootMessage{Error: err})
}

// Write writes message to the client.
func (client *Client) Write(msg *msg_protocol.MsgPackRootMessage) {
	b, err := msgpack.Marshal(msg)
	if err != nil {
		logger.Warn(err.Error())
		return
	}
	client.writer.WriteString(string(b) + "\n")
	client.writer.Flush()
}

// Read reads message from the client.
func (client *Client) Read() {
	for {
		msg, err := client.reader.ReadString('\n')
		if err != nil {
			client.Disconnect()
			logger.Info("stop read goroutine")
			return
		}

		command, args, err := commands.ParseCommand(strings.TrimSpace(msg), " ")
		if err != nil {
			client.WriteError(err)
			continue
		}

		if err := client.Handle(command, args); err != nil {
			client.WriteError(err)
			continue
		}
	}
}

// Disconnect Disconnect
func (client *Client) Disconnect() {
	client.disconnect <- client
	client.conn.Close()
}
