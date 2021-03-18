package tcpconnect

import (
	"errors"
	"strings"

	"github.com/minmax1996/aoimdb/api/commands"
	"github.com/minmax1996/aoimdb/api/msg_protocol"
	db "github.com/minmax1996/aoimdb/internal/aoimdb"
	"github.com/minmax1996/aoimdb/logger"
)

//Handle sends messages to client and return return
func (c *Client) Handle(command commands.Commander, args []string) error {
	switch cmd := command.Name(); {
	case cmd == "auth":
		return c.AuthHandler(args...)
	case cmd == "select" && c.SessionData.authenticated:
		return c.SelectHandler(args...)
	case cmd == "get" && c.SessionData.authenticated:
		return c.GetHandler(args...)
	case cmd == "set" && c.SessionData.authenticated:
		return c.SetHandler(args...)
	case cmd == "exit":
		return c.ExitHandler(args...)
	case !c.SessionData.authenticated:
		return errors.New("not authenticated")
	default:
		return errors.New("unknown command handler")
	}
}

//Send sends command string to establised connection
func (c *Client) AuthHandler(s ...string) error {
	if err := db.DatabaseInstance.AuthentificateByUserPass(s[0], s[1]); err != nil {
		return err
	}

	c.Write(&msg_protocol.MsgPackRootMessage{
		AuthResponse: &msg_protocol.AuthResponse{
			Message: "authenticated",
		},
	})
	c.SessionData.authenticated = true
	return nil
}

//Send sends command string to establised connection
func (c *Client) SelectHandler(s ...string) error {
	c.SessionData.selectedDatabase = s[0]
	db.DatabaseInstance.SelectDatabase(s[0])

	c.Write(&msg_protocol.MsgPackRootMessage{
		SelectResponse: &msg_protocol.SelectResponse{
			SelectedDatabase: s[0],
		},
	})
	return nil
}

//Send sends command string to establised connection
func (c *Client) GetHandler(s ...string) error {
	selectedDatabase := c.SessionData.selectedDatabase
	key := s[0]
	if strings.Contains(s[0], ".") {
		splited := strings.Split(s[0], ".")
		selectedDatabase, key = splited[0], strings.Join(splited[1:], "")
	}

	if len(selectedDatabase) == 0 {
		return errors.New("database not selected")
	}

	db.DatabaseInstance.SelectDatabase(selectedDatabase)

	val, err := db.DatabaseInstance.Get(selectedDatabase, key)
	if err != nil {
		return err
	}

	c.Write(&msg_protocol.MsgPackRootMessage{
		GetResponse: &msg_protocol.GetResponse{
			Key:   key,
			Value: val,
		},
	})
	return nil
}

//Send sends command string to establised connection
func (c *Client) SetHandler(s ...string) error {
	selectedDatabase := c.SessionData.selectedDatabase
	key, value := s[0], s[1]

	if strings.Contains(key, ".") {
		splited := strings.Split(key, ".")
		selectedDatabase, key = splited[0], splited[1]
	}

	if len(selectedDatabase) == 0 {
		return errors.New("database not selected")
	}

	db.DatabaseInstance.SelectDatabase(selectedDatabase)

	err := db.DatabaseInstance.Set(selectedDatabase, key, value)
	if err != nil {
		return err
	}

	c.Write(&msg_protocol.MsgPackRootMessage{
		SetResponse: &msg_protocol.SetResponse{
			Message: "ok",
		},
	})
	return nil
}

//Send sends command string to establised connection
func (c *Client) ExitHandler(s ...string) error {
	c.Write(&msg_protocol.MsgPackRootMessage{
		Message: "Bye",
	})
	c.Disconnect()
	logger.Info("stop read goroutine")
	return nil
}