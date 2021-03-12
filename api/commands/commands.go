package commands

import (
	"fmt"
)

//Commander interface for command from userInput and for server parse
type Commander interface {
	Name() string
	ShowDescription() string
	ShowExample() string
	ValidateUserInput([]string) error
	CallWithArgs(...string) error
}

var commanders []Commander = make([]Commander, 0)

func init() {
	RegisterCommand(NewHelpCommand())
}

//NewAuthCommand constructor standart server commands
func NewAuthCommand(callback func(string, ...string) error) *BaseCommand {
	return NewBaseCommand("auth", 3, 3,
		"auth database command",
		"(Usage: auth user pass)",
		callback)
}

//NewSelectCommand constructor standart server commands
func NewSelectCommand(callback func(string, ...string) error) *BaseCommand {
	return NewBaseCommand("select", 2, 2,
		"select database command",
		"(Usage: select <database_name>)",
		callback)
}

//NewGetCommand constructor standart server commands
func NewGetCommand(callback func(string, ...string) error) *BaseCommand {
	return NewBaseCommand("get", 2, 2,
		"get command",
		"(Usage: get [<databasename>.]<key>)",
		callback)
}

//NewSetCommand constructor standart server commands
func NewSetCommand(callback func(string, ...string) error) *BaseCommand {
	return NewBaseCommand("set", 3, 3,
		"set database command",
		"(Usage: set [<databasename>.]<key> <value>)",
		callback)
}

//NewExitCommand constructor standart server commands
func NewExitCommand(callback func(string, ...string) error) *BaseCommand {
	return NewBaseCommand("exit", 1, 1,
		"exits from server", "",
		callback)
}

//NewHelpCommand constructor standart server commands
func NewHelpCommand() *BaseCommand {
	return NewBaseCommand("help", 1, 2,
		"shows this message", "",
		showAllCommands)
}

//RegisterCommand appends commander object to collection registered commands for help and other
func RegisterCommand(commander Commander) {
	commanders = append(commanders, commander)
}

//GetCommand gets command by name or returns nil if not found
func GetCommand(name string) Commander {
	for _, v := range commanders {
		if v.Name() == name {
			return v
		}
	}
	return nil
}

func showAllCommands(name string, args ...string) error {
	if len(args) == 0 {
		result := "Usage:\n<command> <args>\n\nCommands:\n"
		for _, v := range commanders {
			result += fmt.Sprintf("\t%s\t %s %s\n", v.Name(), v.ShowDescription(), v.ShowExample())
		}
		fmt.Println(result)
	} else if command := GetCommand(args[0]); command != nil {
		fmt.Printf("%s \t %s %s\n", command.Name(), command.ShowDescription(), command.ShowExample())
	}

	return nil
}
