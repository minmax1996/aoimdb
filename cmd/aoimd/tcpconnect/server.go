package tcpconnect

import (
	"net"

	"github.com/minmax1996/aoimdb/api/commands"
	"github.com/minmax1996/aoimdb/logger"
)

//TCPServer struct
type TCPServer struct {
	clients    []*Client
	connect    chan net.Conn
	disconnect chan *Client
	errorChan  chan error
}

// CreateTCPServer creates new server and starts listening for connections.
func CreateTCPServer(errChan chan error) *TCPServer {
	ser := &TCPServer{
		clients:    make([]*Client, 0),
		connect:    make(chan net.Conn),
		disconnect: make(chan *Client),
		errorChan:  errChan,
	}

	commands.RegisterCommand(commands.NewAuthCommand(nil))
	commands.RegisterCommand(commands.NewSelectCommand(nil))
	commands.RegisterCommand(commands.NewGetCommand(nil))
	commands.RegisterCommand(commands.NewSetCommand(nil))
	commands.RegisterCommand(commands.NewKeysCommand(nil))
	commands.RegisterCommand(commands.NewExitCommand(nil))
	commands.RegisterCommand(commands.NewCreateTableCommand(nil))
	commands.RegisterCommand(commands.NewInsertIntoTableCommand(nil))
	commands.RegisterCommand(commands.NewSelectFromTableCommand(nil))

	ser.Listen()

	return ser
}

// Listen listens for connections and messages to broadcast.
func (ser *TCPServer) Listen() {
	go func() {
		for {
			select {
			case conn := <-ser.connect:
				ser.Join(conn)
				logger.InfoFormat("Joined %d", len(ser.clients))
			case dcClien := <-ser.disconnect:
				ser.Remove(dcClien)
				logger.InfoFormat("Disconnected: %d", len(ser.clients))
			}
		}
	}()
}

// Connect passing connection to connect chan to process Join in Listen
func (ser *TCPServer) Connect(conn net.Conn) {
	ser.connect <- conn
}

// Join creates new client and starts listening for client messages.
func (ser *TCPServer) Join(conn net.Conn) {
	client := CreateClient(conn, ser)
	ser.clients = append(ser.clients, client)
}

// Remove disconnected client from server
func (ser *TCPServer) Remove(dcClient *Client) {
	for i, v := range ser.clients {
		if v == dcClient {
			ser.clients = append(ser.clients[:i], ser.clients[i+1:]...) //reslice
			return
		}
	}
}
