package commands

import (
	"errors"
	"fmt"
	"strings"
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

//ParseCommand parse userInput to Commander and validate args, and return them
func ParseCommand(input string, sep string) (Commander, []string, error) {
	uArr := strings.Split(input, sep)
	command := GetCommand(uArr[0])
	if command == nil {
		return nil, nil, errors.New("unknown command `" + uArr[0] + "`")
	}

	if err := command.ValidateUserInput(uArr); err != nil {
		return nil, nil, err
	}

	return command, uArr[1:], nil
}

//NewAuthCommand constructor standart server commands
func NewAuthCommand(callback func(string, ...string) error) Commander {
	return NewBaseCommand("auth", 2, 3,
		"auth database command",
		"(Usage: auth user pass) or (Usage: auth <token>)",
		callback)
}

//NewSelectCommand constructor standart server commands
func NewSelectCommand(callback func(string, ...string) error) Commander {
	return NewBaseCommand("select", 2, 2,
		"select database command",
		"(Usage: select <database_name>)",
		callback)
}

//NewGetCommand constructor standart server commands
func NewGetCommand(callback func(string, ...string) error) Commander {
	return NewBaseCommand("get", 2, 2,
		"get command",
		"(Usage: get [<databasename>.]<key>)",
		callback)
}

//NewSetCommand constructor standart server commands
func NewSetCommand(callback func(string, ...string) error) Commander {
	return NewBaseCommand("set", 3, 3,
		"set database command",
		"(Usage: set [<databasename>.]<key> <value>)",
		callback)
}

//NewCreateTableCommand constructor standart server commands
func NewCreateTableCommand(callback func(string, ...string) error) Commander {
	return NewBaseCommand("tcreate", 3, -1,
		"create table schema in database",
		"(Usage: tcreate <database_name>.tablename name1:int32 name2:string name3:float64)",
		callback)
}

//NewInsertIntoTableCommand constructor standart server commands
func NewInsertIntoTableCommand(callback func(string, ...string) error) Commander {
	return NewBaseCommand("tinsert", 3, -1,
		"insert into table row",
		"(Usage: tinsert <database_name>.tablename 42:myName:0.3123 45:anotherName:3.14",
		callback)
}

//NewSelectFromTableCommand constructor standart server commands
func NewSelectFromTableCommand(callback func(string, ...string) error) Commander {
	return NewBaseCommand("tselect", 3, -1,
		"select fields from table",
		"(Usage: tselect <database_name>.tablename name1 name2 name3",
		callback)
}

//NewSetCommand constructor standart server commands
func NewKeysCommand(callback func(string, ...string) error) Commander {
	return NewBaseCommand("keys", 1, 2,
		"keys database command",
		"(Usage: keys [keysregexp])",
		callback)
}

//NewExitCommand constructor standart server commands
func NewExitCommand(callback func(string, ...string) error) Commander {
	return NewBaseCommand("exit", 1, 1,
		"exits from server", "",
		callback)
}

//NewHelpCommand constructor standart server commands
func NewHelpCommand() Commander {
	return NewBaseCommand("help", 1, 2,
		"shows this message", "",
		showAllCommands)
}

//showAllCommands used for HelpCommand to print all registered commands (or selected one if pass its name to args)
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
