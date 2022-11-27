package tcpconnect

import (
	"errors"
	"strings"

	"github.com/minmax1996/aoimdb/api/commands"
	db "github.com/minmax1996/aoimdb/internal/aoimdb/database"
	"github.com/minmax1996/aoimdb/internal/aoimdb/datatypes"
	"github.com/minmax1996/aoimdb/pkg/logger"
	"github.com/minmax1996/aoimdb/pkg/protocols"
)

// Handle sends messages to client and return return
func (c *Client) Handle(command commands.Commander, args []string) error {
	switch cmd := command.Name(); {
	case cmd == "auth":
		return c.AuthHandler(args...)
	case cmd == "select" && c.SessionData.authenticated:
		return c.SelectHandler(args...)
	case cmd == "keys" && c.SessionData.authenticated:
		return c.KeysHandler(args...)
	case cmd == "get" && c.SessionData.authenticated:
		return c.GetHandler(args...)
	case cmd == "set" && c.SessionData.authenticated:
		return c.SetHandler(args...)
	case cmd == "tcreate" && c.SessionData.authenticated:
		return c.CreateTableHandler(args...)
	case cmd == "tinsert" && c.SessionData.authenticated:
		return c.InsertIntoTableHandler(args...)
	case cmd == "tselect" && c.SessionData.authenticated:
		return c.SelectFromTableHandler(args...)
	case cmd == "exit":
		return c.ExitHandler(args...)
	case !c.SessionData.authenticated:
		return errors.New("not authenticated")
	default:
		return errors.New("unknown command handler")
	}
}

// Send sends command string to establised connection
func (c *Client) AuthHandler(s ...string) error {
	if len(s) == 2 {
		if err := db.AuthentificateByUserPass(s[0], s[1]); err != nil {
			return err
		}
	} else if len(s) == 1 {
		if err := db.AuthentificateByToken(s[0]); err != nil {
			return err
		}
	}

	c.Write(&protocols.Response{
		MessageType: protocols.MessageTypeAuth,
		AuthResponse: &protocols.AuthResponse{
			Message: "authenticated",
		},
	})
	c.SessionData.authenticated = true
	return nil
}

// Send sends command string to establised connection
func (c *Client) SelectHandler(s ...string) error {
	c.SessionData.selectedDatabase = s[0]
	db.SelectDatabase(s[0])

	c.Write(&protocols.Response{
		SelectResponse: &protocols.SelectResponse{
			SelectedDatabase: s[0],
		},
	})
	return nil
}

// Send sends command string to establised connection
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

	db.SelectDatabase(selectedDatabase)

	val, err := db.Get(selectedDatabase, key)
	if err != nil {
		return err
	}

	c.Write(&protocols.Response{
		GetResponse: &protocols.GetResponse{
			Key:   key,
			Value: val,
		},
	})
	return nil
}

// Send sends command string to establised connection
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

	db.SelectDatabase(selectedDatabase)

	err := db.Set(selectedDatabase, key, value)
	if err != nil {
		return err
	}

	c.Write(&protocols.Response{
		MessageType: protocols.MessageTypeSet,
		SetResponse: &protocols.SetResponse{
			Message: "ok",
		},
	})
	return nil
}

// Send sends command string to establised connection
func (c *Client) CreateTableHandler(s ...string) error {
	selectedDatabase := c.SessionData.selectedDatabase
	tableName := s[0]

	if strings.Contains(tableName, ".") {
		splited := strings.Split(tableName, ".")
		selectedDatabase, tableName = splited[0], splited[1]
	}

	if len(selectedDatabase) == 0 {
		return errors.New("database not selected")
	}

	db.SelectDatabase(selectedDatabase)

	var columnNames []string
	var columnTypes []datatypes.Datatype

	for _, v := range s[1:] {
		splitted := strings.Split(v, ":")
		if len(splitted) != 2 {
			continue
		}
		columnNames = append(columnNames, splitted[0])
		columnTypes = append(columnTypes, *datatypes.FromClientString(splitted[1]))
	}

	err := db.CreateTable(selectedDatabase, tableName, columnNames, columnTypes)
	if err != nil {
		return err
	}

	c.Write(&protocols.Response{
		MessageType: protocols.MessageTypeCreateTable,
		CreateTableResponse: &protocols.CreateTableResponse{
			Message: "ok",
		},
	})
	return nil
}

// Send sends command string to establised connection
func (c *Client) InsertIntoTableHandler(s ...string) error {
	selectedDatabase := c.SessionData.selectedDatabase
	tableName := s[0]

	if strings.Contains(tableName, ".") {
		splited := strings.Split(tableName, ".")
		selectedDatabase, tableName = splited[0], splited[1]
	}

	if len(selectedDatabase) == 0 {
		return errors.New("database not selected")
	}

	db.SelectDatabase(selectedDatabase)
	columnNames := strings.Split(s[1], ":")

	//values := make([]interface{}, 0, len(s[2:]))
	for _, v := range s[2:] {
		splitted := strings.Split(v, ":")
		row := make([]interface{}, 0, len(splitted))
		for _, ss := range splitted {
			row = append(row, ss)
		}
		//we can pass interface{} to rows by string
		err := db.InsertIntoTable(selectedDatabase, tableName, columnNames, row)
		if err != nil {
			c.Write(&protocols.Response{
				MessageType: protocols.MessageTypeInsertTable,
				InsertTableResponse: &protocols.InsertTableResponse{
					Message: err.Error(),
				},
			})
			return nil
		}
	}

	c.Write(&protocols.Response{
		MessageType: protocols.MessageTypeInsertTable,
		InsertTableResponse: &protocols.InsertTableResponse{
			Message: "ok",
		},
	})
	return nil
}

// Send sends command string to establised connection
func (c *Client) SelectFromTableHandler(s ...string) error {
	selectedDatabase := c.SessionData.selectedDatabase
	tableName := s[0]

	if strings.Contains(tableName, ".") {
		splited := strings.Split(tableName, ".")
		selectedDatabase, tableName = splited[0], splited[1]
	}

	if len(selectedDatabase) == 0 {
		return errors.New("database not selected")
	}

	db.SelectDatabase(selectedDatabase)

	table := db.SelectFromTable(selectedDatabase, tableName, s[1:])
	if table == nil {
		return nil
	}

	c.Write(&protocols.Response{
		MessageType: protocols.MessageTypeSelectTable,
		SelectTableResponse: &protocols.SelectTableResponse{
			FieldNames: table.ColumnNames,
			Rows:       table.DataRows.Export(),
		},
	})
	return nil
}

// Send sends command string to establised connection
func (c *Client) KeysHandler(s ...string) error {
	//Get client last selected database and no keypatterf by default
	selectedDatabase, keypattern := c.SessionData.selectedDatabase, ""
	if len(s) > 0 {
		//if user provided
		if strings.Contains(s[0], ".") {
			splited := strings.SplitN(s[0], ".", 2)
			selectedDatabase, keypattern = splited[0], splited[1]
		}
	}

	result, err := db.GetKeys(selectedDatabase, keypattern)
	if err != nil {
		return err
	}

	c.Write(&protocols.Response{
		KeysResponse: &protocols.KeysResponse{
			Keys: result,
		},
	})
	return nil
}

// Send sends command string to establised connection
func (c *Client) ExitHandler(s ...string) error {
	c.Write(&protocols.Response{
		Message: "Bye",
	})
	c.Disconnect()
	logger.Info("stop read goroutine")
	return nil
}
