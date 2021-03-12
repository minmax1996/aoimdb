package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/minmax1996/aoimdb/api/commands"
	"github.com/minmax1996/aoimdb/logger"
	"github.com/olekukonko/tablewriter"
)

var (
	username   string
	password   string
	connection net.Conn
	reader     *bufio.Reader
	writer     *bufio.Writer
)

func init() {
	commands.RegisterCommand(commands.NewAuthCommand(Send))
	commands.RegisterCommand(commands.NewSelectCommand(Send))
	commands.RegisterCommand(commands.NewGetCommand(Send))
	commands.RegisterCommand(commands.NewSetCommand(Send))
	commands.RegisterCommand(commands.NewExitCommand(Send))
}

//Send sends command string to establised connection
func Send(name string, s ...string) error {
	writer.WriteString(name + " " + strings.Join(s, " ") + "\n")
	return writer.Flush()
}

func main() {
	var err error

	username = "admin"
	password = "pass"

	//shows help command
	commands.GetCommand("help").CallWithArgs()

	//Open tcp connect to base port
	connection, err = net.Dial("tcp", ":1593")
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	//initiate reader and writer for connection
	reader = bufio.NewReader(connection)
	writer = bufio.NewWriter(connection)

	//try auth by username and password
	//TODO split by args later to do autho auth if userpass (or something else) provided
	commands.GetCommand("auth").CallWithArgs(username, password)

	go startListenResponses()
	commandreader := bufio.NewReader(os.Stdin)
	for {
		command, err := commandreader.ReadString('\n')
		if err != nil {
			break
		}
		//parse userInput to find command and invoke its callback
		if err := parseUserInput(strings.TrimSpace(command)); err != nil {
			logger.Error(err.Error())
		}
	}
	os.Exit(0)
}

func startListenResponses() {
	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			os.Exit(1)
		}
		if strings.HasPrefix(data, "csv>") {
			table, err := tablewriter.NewCSVReader(os.Stdout, csv.NewReader(strings.NewReader(strings.Replace(data, "csv>", "", 1))), false)
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			table.Render()
		} else {
			fmt.Print("< " + data)
		}
	}
}

func parseUserInput(userInput string) error {
	uArr := strings.Split(userInput, " ")
	command := commands.GetCommand(uArr[0])
	if command == nil {
		return errors.New("unknown command")
	}

	if err := command.ValidateUserInput(uArr); err != nil {
		return err
	}
	return command.CallWithArgs(uArr[1:]...)
}
