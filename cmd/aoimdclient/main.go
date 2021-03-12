package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"syscall"

	"github.com/minmax1996/aoimdb/api/commands"
	"github.com/minmax1996/aoimdb/logger"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/crypto/ssh/terminal"
)

var (
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

	//try auth by username and password

	//Open tcp connect to base port
	connection, err = net.Dial("tcp", ":1593")
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	//initiate reader and writer for connection
	reader = bufio.NewReader(connection)
	writer = bufio.NewWriter(connection)
	go startListenResponses()

	if err := handleAuthenticate(); err != nil {
		logger.Error("err in authenticate: " + err.Error())
	}

	//shows help command
	commands.GetCommand("help").CallWithArgs()

	startListenCommands()
	os.Exit(0)
}

func handleAuthenticate() error {
	var username, password string

	flag.StringVar(&username, "u", "", "a string var")
	flag.StringVar(&password, "p", "", "a string var")
	flag.Parse()

	if len(username) > 0 && len(password) > 0 {
		fmt.Println("[Warning] Using a password on the command line interface can be insecure.")
		return commands.GetCommand("auth").CallWithArgs(username, password)
	} else if len(username) > 0 && len(password) == 0 {
		fmt.Print("Enter password: ")
		pass, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		return commands.GetCommand("auth").CallWithArgs(username, string(pass))
	}
	return nil
}

func startListenCommands() {
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
