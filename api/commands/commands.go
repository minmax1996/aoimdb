package commands

import (
	"fmt"
)

//Commander interface for command from userInput and for server parse
type Commander interface {
	Name() string
	ShowDescription() string
	ValidateUserInput([]string) error
	CallWithArgs(...string) error
}

var commanders []Commander = make([]Commander, 0)

func init() {
	RegisterCommand(NewBaseCommand("help", 1, 1,
		"show this message", showAllCommands))
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
	result := ""
	for _, v := range commanders {
		result += v.Name() + "\n" + v.ShowDescription() + "\n"
	}
	fmt.Println(result)
	return nil
}
