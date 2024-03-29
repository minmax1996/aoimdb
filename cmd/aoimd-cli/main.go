package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"

	"github.com/minmax1996/aoimdb/api/commands"
	"github.com/minmax1996/aoimdb/api/msg_protocol"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	connection net.Conn
	reader     *bufio.Reader
	writer     *bufio.Writer
	username   string
	password   string
	host       string
)

func init() {
	commands.RegisterCommand(commands.NewAuthCommand(Send))
	commands.RegisterCommand(commands.NewSelectCommand(Send))
	commands.RegisterCommand(commands.NewGetCommand(Send))
	commands.RegisterCommand(commands.NewSetCommand(Send))
	commands.RegisterCommand(commands.NewCreateTableCommand(Send))
	commands.RegisterCommand(commands.NewInsertIntoTableCommand(Send))
	commands.RegisterCommand(commands.NewSelectFromTableCommand(Send))
	commands.RegisterCommand(commands.NewKeysCommand(Send))
	commands.RegisterCommand(commands.NewExitCommand(Send))

	flag.StringVar(&username, "u", "", "a string var for username")
	flag.StringVar(&password, "p", "", "a string var for password")
	flag.StringVar(&host, "h", "127.0.0.1:1593", "a string var for host to connect")
	flag.Parse()
}

//Send sends command string to establised connection
func Send(name string, s ...string) error {
	writer.WriteString(name + " " + strings.Join(s, " ") + "\n")
	return writer.Flush()
}

func main() {
	var err error

	//Open tcp connect to base port
	connection, err = net.Dial("tcp", host)
	if err != nil {
		fmt.Println("(connect_error) ERR " + err.Error())
		os.Exit(0)
	}
	defer connection.Close()

	//initiate reader and writer for connection
	reader = bufio.NewReader(connection)
	writer = bufio.NewWriter(connection)

	//start goroutine to listen responses from connection
	startListenResponses()

	// try to authenticate user if credentials provided, otherwise user can authenticate by himself later
	if err := handleAuthenticate(); err != nil {
		fmt.Println("(auth_error) ERR " + err.Error())
	}

	//shows help command
	commands.GetCommand("help").CallWithArgs()

	startListenCommands()
	os.Exit(0)
}

func handleAuthenticate() error {
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
	fmt.Println("[Warning] You are not authenticated, please authenticate by typing command 'auth user pass'.")
	fmt.Println("Or you can use flags '-u=user' and promt pass or '-u=user -p=pass' but it cant be insecure.")
	return nil
}

func startListenCommands() {
	commandreader := bufio.NewReader(os.Stdin)
	for {
		userInput, err := commandreader.ReadString('\n')
		if err != nil {
			break
		}

		command, args, err := commands.ParseCommand(strings.TrimSpace(userInput), " ")
		if err != nil {
			fmt.Println("(error) ERR " + err.Error())
			continue
		}

		if err := command.CallWithArgs(args...); err != nil {
			fmt.Println("(error) ERR " + err.Error())
			continue
		}
	}
}

func startListenResponses() {
	go func() {
		for {
			data, err := reader.ReadString('\n')
			if err != nil {
				os.Exit(1)
			}

			var item msg_protocol.MsgPackRootMessage
			err = msgpack.Unmarshal([]byte(data), &item)
			if err != nil {
				fmt.Println("(error) ERR " + err.Error())
				continue
			}

			if item.Error != nil {
				fmt.Println("(error) ERR " + item.Error.Error())
				continue
			}

			switch {
			case item.AuthResponse != nil:
				fmt.Println(*item.AuthResponse)
			case item.SelectResponse != nil:
				fmt.Println(*item.SelectResponse)
			case item.GetResponse != nil:
				fmt.Println(*item.GetResponse)
			case item.SetResponse != nil:
				fmt.Println(*item.SetResponse)
			case item.KeysResponse != nil:
				fmt.Println(*item.KeysResponse)
			case item.CreateTableResponse != nil:
				fmt.Println(item.CreateTableResponse.Message)
			case item.InsertTableResponse != nil:
				fmt.Println(item.InsertTableResponse.Message)
			case item.SelectTableResponse != nil:
				fmt.Println(*item.SelectTableResponse)
			default:
				fmt.Println(item.Message)
			}

			// if strings.HasPrefix(item.Message, "csv>") {
			// 	table, err := tablewriter.NewCSVReader(os.Stdout, csv.NewReader(strings.NewReader(strings.Replace(item.Message, "csv>", "", 1))), false)
			// 	if err != nil {
			// 		logger.Error(err.Error())
			// 		continue
			// 	}
			// 	table.Render()
			// } else {
			// 	fmt.Print("< " + item.Message)
			// }
		}
	}()
}
